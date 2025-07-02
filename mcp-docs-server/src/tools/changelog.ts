/**
 * Get recent changes, updates, and new features in Vapi
 */
export async function getChangelog(
  version?: string,
  limit: number = 10,
  type: string = "all"
): Promise<string> {
  try {
    // Sample changelog data - in a real implementation, this would come from a changelog API or file
    const changelogEntries = [
      {
        version: "1.8.0",
        date: "2024-01-25",
        type: "features",
        title: "Enhanced Voice Settings & Deepgram Aura Support",
        changes: [
          "Added support for Deepgram Aura voice models",
          "New voice interruption settings for better conversation flow",
          "Enhanced background noise detection and cancellation",
          "Support for custom voice cloning with ElevenLabs"
        ]
      },
      {
        version: "1.7.5",
        date: "2024-01-22",
        type: "fixes",
        title: "Bug Fixes & Performance Improvements",
        changes: [
          "Fixed webhook delivery reliability issues",
          "Improved call connection stability",
          "Resolved audio quality issues on mobile networks",
          "Fixed function calling timeout errors"
        ]
      },
      {
        version: "1.7.0",
        date: "2024-01-18",
        type: "features",
        title: "Advanced Function Calling & Tool Integration",
        changes: [
          "New parallel function calling support",
          "Enhanced tool parameter validation",
          "Webhook retry mechanism with exponential backoff",
          "Support for streaming responses from functions"
        ]
      },
      {
        version: "1.6.8",
        date: "2024-01-15",
        type: "fixes",
        title: "Call Management Improvements",
        changes: [
          "Fixed call status tracking inconsistencies",
          "Improved error messages for failed calls",
          "Enhanced call recording quality",
          "Fixed timezone issues in call logs"
        ]
      },
      {
        version: "1.6.5",
        date: "2024-01-12",
        type: "features",
        title: "Multi-Language Support & Localization",
        changes: [
          "Added support for 15+ languages",
          "New language detection for incoming calls",
          "Enhanced accent handling for voice recognition",
          "Localized error messages and prompts"
        ]
      },
      {
        version: "1.6.0",
        date: "2024-01-08",
        type: "breaking",
        title: "API V2 Release",
        changes: [
          "New RESTful API structure with improved consistency",
          "Updated authentication flow with refresh tokens",
          "Enhanced error handling with detailed error codes",
          "Migration guide available for V1 users"
        ]
      },
      {
        version: "1.5.12",
        date: "2024-01-05",
        type: "features",
        title: "Dashboard Enhancements",
        changes: [
          "New real-time call monitoring dashboard",
          "Enhanced analytics with custom date ranges",
          "Bulk operations for assistant management",
          "Export functionality for call logs and analytics"
        ]
      },
      {
        version: "1.5.8",
        date: "2024-01-02",
        type: "fixes",
        title: "Holiday Bug Fixes",
        changes: [
          "Fixed assistant configuration validation",
          "Resolved phone number formatting issues",
          "Improved error handling for invalid API keys",
          "Fixed memory leaks in long-running calls"
        ]
      }
    ];

    // Filter by version if specified
    let filteredEntries = changelogEntries;
    if (version) {
      filteredEntries = changelogEntries.filter(entry => 
        entry.version === version || 
        entry.version.includes(version)
      );
    }

    // Filter by type if specified
    if (type !== "all") {
      filteredEntries = filteredEntries.filter(entry => entry.type === type);
    }

    // Limit results
    filteredEntries = filteredEntries.slice(0, limit);

    if (filteredEntries.length === 0) {
      return `# 📋 No Changelog Entries Found

${version ? `No changelog entries found for version "${version}".` : "No changelog entries found."}

## Available Versions:
${changelogEntries.map(entry => `- **v${entry.version}** (${entry.date}) - ${entry.title}`).join('\n')}

## Filter Options:
- **Type:** features, fixes, breaking, all
- **Version:** Specify exact version (e.g., "1.8.0")

Try searching without filters or check https://docs.vapi.ai/changelog`;
    }

    let response = `# 📋 Vapi Changelog\n\n`;
    
    if (version) {
      response += `Showing changes for version "${version}"\n`;
    } else {
      response += `Showing latest ${filteredEntries.length} changelog entries\n`;
    }
    
    if (type !== "all") {
      response += `**Filter:** ${type}\n`;
    }
    response += `\n`;

    filteredEntries.forEach((entry, index) => {
      // Add type emoji
      let typeEmoji = "📝";
      switch (entry.type) {
        case "features": typeEmoji = "✨"; break;
        case "fixes": typeEmoji = "🐛"; break;
        case "breaking": typeEmoji = "⚠️"; break;
      }

      response += `## ${typeEmoji} v${entry.version} - ${entry.title}\n`;
      response += `**Released:** ${entry.date} | **Type:** ${entry.type}\n\n`;
      
      entry.changes.forEach(change => {
        response += `- ${change}\n`;
      });
      
      response += `\n`;
      
      if (index < filteredEntries.length - 1) {
        response += "---\n\n";
      }
    });

    response += `\n## 🔗 Additional Resources\n\n`;
    response += `- **[Complete Changelog](https://docs.vapi.ai/changelog)** - Full version history\n`;
    response += `- **[Migration Guides](https://docs.vapi.ai/migrations)** - Upgrade instructions\n`;
    response += `- **[Breaking Changes](https://docs.vapi.ai/breaking-changes)** - Important updates\n`;
    response += `- **[Release Notes](https://github.com/VapiAI/releases)** - GitHub releases\n`;
    response += `- Use \`search_documentation\` to find related guides`;

    return response;
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `❌ Error fetching changelog: ${errorMessage}\n\nPlease try again or visit https://docs.vapi.ai/changelog`;
  }
} 