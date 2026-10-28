import hashlib,io,os,pathlib,shutil,subprocess,tarfile,tempfile
root=pathlib.Path(__file__).resolve().parents[2];os.chdir(root)
version='0.2.0-alpha.1';name=f'spotify-tui_{version}_linux_amd64';archive=name+'.tar.gz'
base=pathlib.Path(tempfile.mkdtemp(prefix='release test ',dir=root/'work'))
fake=base/'fakebin';fake.mkdir();gh=fake/'gh'
gh.write_text('''#!/usr/bin/env bash
set -e
if [[ "$1" == auth ]]; then exit 0; fi
while [[ $# -gt 0 ]]; do
 if [[ "$1" == --dir ]]; then out=$2;shift 2;else shift;fi
done
cp "$FIXTURE_DIR/$FIXTURE_ARCHIVE" "$out/$FIXTURE_ARCHIVE"
cp "$FIXTURE_DIR/SHA256SUMS" "$out/SHA256SUMS"
''');gh.chmod(0o755)
env=os.environ.copy();env['PATH']=str(fake)+os.pathsep+env['PATH'];env['FIXTURE_ARCHIVE']=archive
for key in ['XDG_CONFIG_HOME','XDG_DATA_HOME','XDG_CACHE_HOME','XDG_STATE_HOME','XDG_RUNTIME_DIR']:env[key]=str(base/key)
fixture=base/'fixtures';fixture.mkdir();env['FIXTURE_DIR']=str(fixture)
def hashes():
 (fixture/'SHA256SUMS').write_text(hashlib.sha256((fixture/archive).read_bytes()).hexdigest()+'  '+archive+'\n')
def run(case,success):
 destination=base/case
 p=subprocess.run(['bash','scripts/download-linux.sh','--destination',str(destination)],env=env,capture_output=True,text=True)
 assert (p.returncode==0)==success,(case,p.stdout,p.stderr)
 assert destination.exists()==success,case
 assert not list(base.glob('.spotify-tui-linux.*')),case+' temporary file leak'
 return destination
shutil.copyfile(root/'dist/release'/archive,fixture/archive);hashes()
valid=run('valid download',True)
assert (valid/'spotify-tui').read_bytes()==(root/'dist/spotify-tui').read_bytes()
with (fixture/archive).open('ab') as f:f.write(b'corrupt')
run('bad checksum',False)
for kind in ['traversal','symlink']:
 with tarfile.open(fixture/archive,'w:gz') as t:
  member=tarfile.TarInfo(name+'/../escaped' if kind=='traversal' else name+'/link')
  if kind=='symlink':member.type=tarfile.SYMTYPE;member.linkname='/tmp/outside'
  t.addfile(member,io.BytesIO(b''))
 hashes();run(kind,False)
assert not (base/'escaped').exists()
existing=base/'existing';existing.mkdir();(existing/'keep').write_text('keep')
p=subprocess.run(['bash','scripts/download-linux.sh','--destination',str(existing)],env=env,capture_output=True)
assert p.returncode!=0 and (existing/'keep').read_text()=='keep'
prefix=base/'install prefix'
subprocess.run(['bash',str(valid/'install.sh'),'--prefix',str(prefix)],check=True,capture_output=True)
assert subprocess.check_output([str(prefix/'bin/spotify-tui'),'version'],text=True,env=env).strip()=='spotify-tui 0.2.0-alpha.1'
assert subprocess.check_output([str(prefix/'bin/spotify-tui-playback'),'--version'],text=True,env=env).strip()=='spotify-tui-playback 0.2.0-alpha.1'
subprocess.run(['bash',str(valid/'install.sh'),'--prefix',str(prefix)],check=True,capture_output=True)
print('Linux archive download, checksum failure, path/link rejection, overwrite refusal, temp cleanup, package installation/upgrade and version pairing (paths with spaces): PASS')
