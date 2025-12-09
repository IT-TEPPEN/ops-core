/**
 * Substitutes variables in a content string with provided values
 * @param content - The content string with variable placeholders in {{variable_name}} format
 * @param values - An object mapping variable names to their values
 * @returns The content with all variable placeholders replaced
 */
export function substituteVariables(
  content: string,
  values: { [key: string]: any }
): string {
  let result = content;

  for (const [name, value] of Object.entries(values)) {
    // Create a regex to match {{variable_name}} (case-sensitive, exact match)
    const regex = new RegExp(`\\{\\{${escapeRegExp(name)}\\}\\}`, 'g');
    
    // Convert value to string, handling different types
    const strValue = String(value ?? '');
    
    // Replace all occurrences
    result = result.replace(regex, strValue);
  }

  return result;
}

/**
 * Escapes special regex characters in a string
 * @param str - The string to escape
 * @returns The escaped string safe for use in regex
 */
function escapeRegExp(str: string): string {
  return str.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
}
