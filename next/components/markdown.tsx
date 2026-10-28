import ReactMarkdown from "react-markdown";
import remarkGfm from "remark-gfm";
import path from "node:path";
import { repo } from "@/lib/releases";
export function Markdown({
  children,
  base = "",
}: {
  children: string;
  base?: string;
}) {
  const content = children.replace(/<a id="[^"]+"><\/a>\s*/g, "");
  return (
    <div className="prose">
      <ReactMarkdown
        remarkPlugins={[remarkGfm]}
        components={{
          h2: ({ children }) => <h2 data-reveal>{children}</h2>,
          h3: ({ children }) => (
            <h3 data-reveal
              id={String(children)
                .toLowerCase()
                .replace(/[^a-z0-9]/g, "")}
            >
              {children}
            </h3>
          ),
          a: ({ href, children }) => {
            const match = href?.match(/(?:^|\/)(v\d+\.\d+\.\d+)\.md/);
            const local =
              href && !/^(https?:|#|\/)/.test(href)
                ? path.posix.normalize(path.posix.join(base, href))
                : null;
            const target = match
              ? `/changelog/${match[1]}${href?.includes("#") ? "#" + href.split("#")[1] : ""}`
              : local?.startsWith("docs/MEMORY.md")
                ? "/memory"
                : local?.startsWith("docs/SETUP.md")
                ? "/docs"
                : local?.startsWith("CHANGELOG.md")
                  ? `/changelog${local.includes("#") ? "#" + local.split("#")[1] : ""}`
                  : local
                    ? `${repo}/blob/main/${local}`
                    : href;
            return <a href={target}>{children}</a>;
          },
        }}
      >
        {content}
      </ReactMarkdown>
    </div>
  );
}
