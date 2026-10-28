#!/usr/bin/env python3
"""Native release tests on Linux and macOS, using isolated installations/data."""
import hashlib
import importlib.util
import io
import os
from pathlib import Path
import shutil
import subprocess
import tarfile
import tempfile
import unittest
from unittest.mock import patch

ROOT = Path(__file__).resolve().parents[2]
spec = importlib.util.spec_from_file_location("release", ROOT / "scripts/release.py")
release = importlib.util.module_from_spec(spec)
spec.loader.exec_module(release)
VERSION = release.VERSION
NAME = "spotify-tui_" + VERSION + "_" + release.target()


class ReleaseTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix="spotify release test ")
        self.addCleanup(self.temp.cleanup)
        self.base = Path(self.temp.name)
        self.archive = ROOT / "dist/release" / (NAME + ".tar.gz")
        self.source = self.base / "unpacked"
        self.source.mkdir()
        release.extract(self.archive, self.source, NAME)

    def test_native_pair_install_reinstall_and_unrelated_file_refusal(self):
        prefix = self.base / "prefix with spaces"
        release.install(self.source, prefix)
        old = (prefix / "lib/spotify-tui/current").resolve()
        release.install(self.source, prefix)
        new = (prefix / "lib/spotify-tui/current").resolve()
        self.assertNotEqual(old, new)
        self.assertTrue(old.is_dir())
        env = dict(os.environ, HOME=str(self.base / "home"))
        for name in release.BINARIES:
            flag = "version" if name == "spotify-tui" else "--version"
            self.assertEqual(subprocess.check_output([str(prefix / "bin" / name), flag], env=env, text=True).strip(), name + " " + VERSION)
            self.assertEqual((prefix / "bin" / name).resolve().parent, new)
        unrelated = self.base / "unrelated/bin"
        unrelated.mkdir(parents=True)
        (unrelated / "spotify-tui").write_text("keep me")
        with self.assertRaises(ValueError):
            release.install(self.source, unrelated.parent)
        self.assertEqual((unrelated / "spotify-tui").read_text(), "keep me")

    def test_checksum_download_cleanup_and_existing_destination(self):
        def fetch(url, path, limit=release.LIMIT):
            source = self.archive if url.endswith(".tar.gz") else self.archive.parent / "SHA256SUMS"
            shutil.copy2(source, path)
        destination = self.base / "download with spaces"
        with patch.object(release, "fetch", side_effect=fetch):
            release.download(VERSION, destination)
            self.assertEqual((destination / "spotify-tui").read_bytes(), (self.source / "spotify-tui").read_bytes())
            with self.assertRaises(ValueError):
                release.download(VERSION, destination)
        self.assertFalse(list(self.base.glob(".spotify-download-*")))
        bad = self.base / "bad.tar.gz"
        bad.write_bytes(b"corrupt")
        sums = self.base / "SHA256SUMS"
        sums.write_text("0" * 64 + "  bad.tar.gz\n")
        with self.assertRaises(ValueError):
            release.verify_checksum(bad, sums)
        digest = hashlib.sha256(bad.read_bytes()).hexdigest()
        sums.write_text((digest + "  bad.tar.gz\n") * 2)
        with self.assertRaises(ValueError):
            release.verify_checksum(bad, sums)

    def test_unsafe_archive_members_refused_before_writing(self):
        for case in ("traversal", "link", "duplicate", "wrong-root"):
            archive = self.base / (case + ".tar.gz")
            with tarfile.open(archive, "w:gz") as output:
                name = NAME + "/file"
                if case == "traversal": name = NAME + "/../escaped"
                if case == "wrong-root": name = "elsewhere/file"
                member = tarfile.TarInfo(name)
                if case == "link":
                    member.type = tarfile.SYMTYPE
                    member.linkname = "/tmp/outside"
                output.addfile(member, io.BytesIO())
                if case == "duplicate": output.addfile(member, io.BytesIO())
            destination = self.base / case
            destination.mkdir()
            with self.assertRaises(ValueError):
                release.extract(archive, destination, NAME)
            self.assertEqual(list(destination.iterdir()), [])

    def test_platform_mapping_and_real_installer_wrapper(self):
        with patch.object(release.platform, "system", return_value="Darwin"), patch.object(release.platform, "machine", return_value="arm64"):
            self.assertEqual(release.target(), "darwin_arm64")
        with patch.object(release.platform, "system", return_value="Darwin"), patch.object(release.platform, "machine", return_value="x86_64"):
            with self.assertRaises(ValueError): release.target()
        subprocess.run(["bash", str(self.source / "install.sh"), "--prefix", str(self.base / "wrapper prefix")], check=True)


if __name__ == "__main__":
    unittest.main()
