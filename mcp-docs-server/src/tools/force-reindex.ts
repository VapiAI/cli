import { DocsFetcher } from "../utils/docs-fetcher.js";

const docsFetcher = new DocsFetcher();

/**
 * Force a complete re-index of documentation for testing purposes
 */
export async function forceReindex(
  clearCache: boolean = true,
  skipVectorIndex: boolean = false
): Promise<string> {
  try {
    let response = `# 🔄 Force Reindex Started\n\n`;
    
    response += `**Timestamp:** ${new Date().toISOString()}\n`;
    response += `**Clear Cache:** ${clearCache ? 'Yes' : 'No'}\n`;
    response += `**Skip Vector Index:** ${skipVectorIndex ? 'Yes' : 'No'}\n\n`;
    
    if (clearCache) {
      response += `## 🗑️ Clearing Cache\n\n`;
      await docsFetcher.invalidateCache();
      response += `✅ Cache cleared successfully\n\n`;
    }
    
    response += `## 📥 Fetching Fresh Documentation\n\n`;
    
    // Force fresh fetch of documentation
    const startTime = Date.now();
    const docs = await docsFetcher.getDocumentationStructure();
    const fetchTime = Date.now() - startTime;
    
    response += `✅ Documentation fetched: ${docs.pages.length} pages in ${fetchTime}ms\n\n`;
    
    if (!skipVectorIndex) {
      response += `## 🧠 Re-indexing Vector Search\n\n`;
      
      const vectorStartTime = Date.now();
      const examples = await docsFetcher.getExamples();
      const vectorTime = Date.now() - vectorStartTime;
      
      response += `✅ Vector index rebuilt: ${examples.length} examples in ${vectorTime}ms\n\n`;
      
      // Get vector index stats
      const indexSize = docsFetcher.getVectorIndexSize();
      response += `📊 **Vector Index Stats:**\n`;
      response += `- Total indexed documents: ${indexSize}\n`;
      response += `- Vector search model: Available\n\n`;
    }
    
    response += `## 📊 Final Statistics\n\n`;
    response += `- **Total documentation pages:** ${docs.pages.length}\n`;
    response += `- **Total processing time:** ${Date.now() - startTime}ms\n`;
    response += `- **Cache status:** Fresh\n`;
    response += `- **Vector search:** ${skipVectorIndex ? 'Skipped' : 'Rebuilt'}\n\n`;
    
    response += `## 🎯 Testing Suggestions\n\n`;
    response += `- Use \`search_documentation\` to test semantic search\n`;
    response += `- Use \`get_examples\` to verify example extraction\n`;
    response += `- Try queries like "MCP server" or "phone calls"\n`;
    response += `- Check for recent documentation updates\n\n`;
    
    response += `## ✅ Reindex Complete\n\n`;
    response += `The documentation has been freshly fetched and re-indexed. All tools should now have the latest content available.`;
    
    return response;
    
  } catch (error) {
    const errorMessage = error instanceof Error ? error.message : "Unknown error";
    return `# ❌ Reindex Failed

Failed to reindex documentation: ${errorMessage}

## 🛠️ Troubleshooting:
- Check internet connectivity for repository access
- Verify git is installed and accessible
- Ensure sufficient disk space for caching
- Try running with skipVectorIndex=true if vector indexing fails

## 🔄 Retry Options:
- Try \`force_reindex\` with \`clearCache: false\` to keep existing cache
- Try \`force_reindex\` with \`skipVectorIndex: true\` to skip vector rebuilding
- Check the main documentation tools are still functional

## 📋 Manual Fallback:
If reindexing continues to fail, the existing cached documentation should still be available for search.`;
  }
} 