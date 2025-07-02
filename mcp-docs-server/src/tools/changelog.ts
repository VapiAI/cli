import { getChangelog as getChangelogDocs } from "../utils/documentation-data";

/**
 * Get recent changes, updates, and new features in Vapi
 */
export async function getChangelog(
  version?: string,
  limit: number = 10,
  type: string = "all"
): Promise<string> {
  try {
    // Get changelog from real Vapi documentation
    const changelogEntries = await getChangelogDocs(limit + 5); // Get extra to allow for filtering

    // Filter by version if specified
    let filteredEntries = changelogEntries;
    if (version) {
      filteredEntries = changelogEntries.filter(entry => 
        entry.title.includes(version) || 
        entry.url.includes(version)
      );
    }

    // Note: type filtering not applicable to real docs structure, but we keep the parameter for compatibility

    // Limit results
    filteredEntries = filteredEntries.slice(0, limit);

    if (filteredEntries.length === 0) {
      return `# 📋 No Changelog Entries Found

${version ? `No changelog entries found for version "${version}".` : "No changelog entries found."}

## Recent Updates Available:
${changelogEntries.slice(0, 10).map(entry => `- **${entry.title}** - ${entry.description}`).join('\n')}

📚 **Complete Changelog:** https://docs.vapi.ai/changelog

Try searching without filters or visit the documentation directly.`;
    }

    let response = `# 📋 Vapi Changelog\n\n`;
    
    if (version) {
      response += `Showing changes for "${version}"\n`;
    } else {
      response += `Showing latest ${filteredEntries.length} changelog entries\n`;
    }
    response += `\n`;

    filteredEntries.forEach((entry, index) => {
      // Extract date from URL if possible (e.g., /changelog/2025/1/15.mdx)
      const dateMatch = entry.url.match(/\/changelog\/(\d{4})\/(\d{1,2})\/(\d{1,2})/);
      let dateStr = "";
      if (dateMatch && dateMatch[1] && dateMatch[2] && dateMatch[3]) {
        const year = dateMatch[1];
        const month = dateMatch[2].padStart(2, '0');
        const day = dateMatch[3].padStart(2, '0');
        dateStr = `${year}-${month}-${day}`;
      }

      response += `## 📝 ${entry.title}\n`;
      if (dateStr) {
        response += `**Released:** ${dateStr}\n`;
      }
      response += `**Category:** ${entry.category}\n\n`;
      response += `${entry.description}\n\n`;
      response += `**📖 View Details:** ${entry.url}\n\n`;
      
      if (index < filteredEntries.length - 1) {
        response += "---\n\n";
      }
    });

    response += `\n## 🔗 Additional Resources\n\n`;
    response += `- **[Complete Changelog](https://docs.vapi.ai/changelog)** - Full version history\n`;
    response += `- **[Documentation](https://docs.vapi.ai)** - Latest features and guides\n`;
    response += `- **[API Reference](https://docs.vapi.ai/api-reference)** - API updates\n`;
    response += `- Use \`search_documentation\` to find information about specific features\n\n`;
    
    response += `💡 **Tip:** Visit the changelog links above for detailed release notes, migration guides, and breaking change information.`;

    return response;
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `❌ Error fetching changelog: ${errorMessage}\n\nPlease visit https://docs.vapi.ai/changelog for the latest updates`;
  }
} 