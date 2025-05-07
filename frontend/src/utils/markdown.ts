import yaml from "js-yaml";
import { z } from "zod";

const schemaMetadata = z.object({
  title: z.string().optional(),
  description: z.string().optional(),
  owner: z.string().optional(),
  version: z.string().optional(),
  type: z.string().optional(),
  tags: z.array(z.string()).optional(),
});

export type TSchemaMetadata = z.infer<typeof schemaMetadata>;

export function parseFrontmatter(markdown: string): {
  content: string;
  meta: TSchemaMetadata;
} {
  const result: { content: string; meta: TSchemaMetadata } = {
    content: markdown,
    meta: {},
  };

  // Check if the content starts with frontmatter delimiters (---)
  if (!markdown.startsWith("---")) {
    return result;
  }

  // Find the end of the frontmatter block
  const endOfFrontmatter = markdown.indexOf("---", 3);
  if (endOfFrontmatter === -1) {
    return result;
  }

  // Extract frontmatter text
  const yamlText = markdown.substring(3, endOfFrontmatter).trim();

  // Parse the frontmatter using YAML
  try {
    result.meta = schemaMetadata.parse(yaml.load(yamlText));
  } catch (e) {
    console.error("Failed to parse frontmatter:", e);
    result.meta = {};
  }

  // Return content without frontmatter
  result.content = markdown.substring(endOfFrontmatter + 3).trim();
  return result;
}
