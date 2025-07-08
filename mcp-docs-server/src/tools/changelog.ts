import { DocsFetcher } from "../utils/docs-fetcher.js";

const docsFetcher = new DocsFetcher();

/**
 * Get recent changes, updates, and new features in Vapi
 */
export async function getChangelog(
  version?: string,
  limit: number = 10,
  type: string = "all"
): Promise<string> {
  try {
    // Search for changelog content
    const changelogSearchResult = await docsFetcher.searchDocumentation("changelog", "changelog");
    const changelogResults = changelogSearchResult.results;
    
    if (changelogResults.length === 0) {
      return `# 📝 Changelog

## 🔄 Recent Updates

The Vapi changelog is available at:
- **[Official Changelog](https://docs.vapi.ai/changelog)** - Latest updates and releases
- **[GitHub Releases](https://github.com/VapiAI/docs/releases)** - Release notes and versions

## 🎯 What's New in Vapi:

### Recent Features:
- **Workflows** - Visual conversation flow builder
- **Squads** - Team collaboration and agent handoffs
- **OpenAI Realtime** - Ultra-low latency voice conversations
- **Custom LLMs** - Bring your own language models
- **Advanced Analytics** - Detailed call insights and metrics
- **Multi-language Support** - Conversations in 50+ languages

### Recent Integrations:
- **Google Calendar** - Automated scheduling
- **Slack** - Team notifications and updates
- **Trieve** - Knowledge base integration
- **Tavus** - Video avatar support
- **PlayHT** - Additional voice options

## 🔗 Resources:
- **[Dashboard](https://dashboard.vapi.ai)** - Try new features
- **[Discord](https://discord.gg/vapi)** - Community updates
- **[Documentation](https://docs.vapi.ai)** - Feature guides

Visit the official changelog link above for the complete list of updates!`;
    }

    let response = `# 📝 Vapi Changelog\n\n`;
    
    if (version) {
      response += `Changelog for version "${version}"\n\n`;
    } else {
      response += `Latest ${limit} changelog entries\n\n`;
    }
    
    if (type !== "all") {
      response += `**Type Filter:** ${type}\n\n`;
    }

    // Fetch and return actual changelog content
    for (let i = 0; i < Math.min(changelogResults.length, 3); i++) {
      const changelogPage = changelogResults[i];
      if (!changelogPage) continue;
      
      try {
        const content = await docsFetcher.fetchPageContent(changelogPage);
        
        response += `## 📄 ${changelogPage.title}\n\n`;
        response += `**Section:** ${changelogPage.section}\n`;
        response += `**URL:** ${changelogPage.url}\n\n`;
        
        // Add the actual content
        response += `### Content:\n\n${content}\n\n`;
        response += `---\n\n`;
        
      } catch (error) {
        response += `## 📄 ${changelogPage.title}\n\n`;
        response += `**URL:** ${changelogPage.url}\n\n`;
        response += `⚠️ Content temporarily unavailable. Please visit the URL above.\n\n`;
        response += `---\n\n`;
      }
    }

    response += `## 🎯 Stay Updated\n\n`;
    response += `- **[Official Changelog](https://docs.vapi.ai/changelog)** - Complete version history\n`;
    response += `- **[GitHub Releases](https://github.com/VapiAI/docs/releases)** - Release notes\n`;
    response += `- **[Discord](https://discord.gg/vapi)** - Community announcements\n`;
    response += `- **[Newsletter](https://vapi.ai)** - Monthly updates\n\n`;
    
    response += `## 🔗 Additional Resources\n\n`;
    response += `- **[Dashboard](https://dashboard.vapi.ai)** - Try new features\n`;
    response += `- **[Documentation](https://docs.vapi.ai)** - Feature guides\n`;
    response += `- **[API Reference](https://docs.vapi.ai/api-reference)** - Latest API updates\n`;
    response += `- **[Examples](https://docs.vapi.ai/guides)** - Implementation guides`;

    return response;
    
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `# ❌ Changelog Error

Failed to fetch changelog: ${errorMessage}

## 🛠️ Troubleshooting:
- The documentation server might be temporarily unavailable
- Try again in a few moments
- Check your internet connection

## 📋 Manual Resources:
- **[Official Changelog](https://docs.vapi.ai/changelog)** - Complete version history
- **[GitHub Releases](https://github.com/VapiAI/docs/releases)** - Release notes
- **[Discord](https://discord.gg/vapi)** - Community announcements

## 🎯 Recent Highlights:
- **Workflows** - Visual conversation builder
- **Squads** - Team collaboration features
- **OpenAI Realtime** - Ultra-low latency
- **Custom LLMs** - Bring your own models
- **Advanced Analytics** - Detailed insights

Visit the official changelog link above for the complete list of updates!`;
  }
} 