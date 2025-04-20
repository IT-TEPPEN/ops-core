import ReactMarkdown from "react-markdown";

const markdownContent = `
# My First Blog Post

This is my **first blog post** written in *Markdown*.

## Introduction

Markdown is a lightweight markup language with plain-text-formatting syntax.

- Item 1
- Item 2
  - Nested Item

\`\`\`javascript
function greet(name) {
  console.log(\`Hello, \${name}!\`);
}

greet('World');
\`\`\`

> This is a blockquote.

You can find the source code on [GitHub](https://github.com).

---

Hope you enjoy reading!

## Conclusion

gathering information from various sources is essential for learning.
`;

// Ensure this component is exported as default
function BlogPage() {
  return (
    // Apply Tailwind Typography styles and some basic layout/theming
    <div className="prose lg:prose-xl dark:prose-invert mx-auto p-6 md:p-8 bg-white dark:bg-gray-800 rounded-lg shadow-md my-8 max-w-3xl">
      <article className="prose lg:prose-xl">
        <ReactMarkdown>{markdownContent}</ReactMarkdown>
      </article>
    </div>
  );
}

export default BlogPage; // Make sure this line exists and is correct
