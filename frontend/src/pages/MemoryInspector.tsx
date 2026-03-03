import { useState } from 'react';
import { useMemories, useSearchMemories, useGiveMemoryFeedback } from '../hooks/useMemories';
import { Memory } from '../types';

const MemoryInspector: React.FC = () => {
  const [searchQuery, setSearchQuery] = useState('');
  const [selectedMemory, setSelectedMemory] = useState<Memory | null>(null);
  
  const { data: allMemories, isLoading: listLoading } = useMemories();
  const { data: searchResults, isFetching: searchLoading } = useSearchMemories(searchQuery);
  const feedbackMutation = useGiveMemoryFeedback();

  const memories = searchQuery.length > 2 ? searchResults : allMemories;
  const isLoading = listLoading || (searchQuery.length > 2 && searchLoading);

  const handleFeedback = async (id: string, confidence: number) => {
    try {
      await feedbackMutation.mutateAsync({ id, feedback: { confidence } });
    } catch (err) {
      console.error("Feedback failed:", err);
    }
  };

  const getConfidenceColor = (score: number) => {
    if (score > 0.8) return 'bg-green-500';
    if (score > 0.5) return 'bg-yellow-500';
    return 'bg-red-500';
  };

  return (
    <div className="flex h-full w-full overflow-hidden bg-bg-base">
      <div className="flex-1 flex flex-col min-w-0 overflow-hidden">
        <div className="flex-none p-6 border-b border-border-base bg-bg-surface space-y-4 shadow-sm z-10">
          <div className="flex justify-between items-center">
            <div>
              <h1 className="text-2xl font-bold text-text-base">Memory Inspector</h1>
              <p className="text-sm text-text-muted">Audit and evolve semantic memories across the brain.</p>
            </div>
            <div className="flex items-center space-x-2">
              <label htmlFor="source-filter" className="text-xs font-medium text-text-muted uppercase tracking-wider">Source:</label>
              <select 
                id="source-filter"
                name="source-filter"
                className="bg-bg-base border border-border-base rounded-md px-2 py-1 text-xs outline-none focus:ring-1 focus:ring-primary"
              >
                <option>All Sources</option>
                <option>MCP Tool</option>
                <option>Workflow</option>
              </select>
            </div>
          </div>

          <div className="relative">
            <span className="absolute inset-y-0 left-0 pl-3 flex items-center text-gray-400">
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </span>
            <input 
              id="memory-search"
              name="memory-search"
              type="text"
              placeholder="Search memories semantically..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full bg-bg-base border border-border-base rounded-xl pl-10 pr-4 py-3 focus:ring-2 focus:ring-primary outline-none transition-all"
            />
            {searchLoading && (
              <div className="absolute inset-y-0 right-0 pr-3 flex items-center">
                <div className="animate-spin h-4 w-4 border-2 border-primary border-t-transparent rounded-full"></div>
              </div>
            )}
          </div>
        </div>

        <div className="flex-1 overflow-y-auto p-6">
          {isLoading ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {[...Array(6)].map((_, i) => (
                <div key={i} className="h-32 bg-bg-surface border border-border-base rounded-xl animate-pulse"></div>
              ))}
            </div>
          ) : (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {memories?.map(memory => (
                <button
                  key={memory.id}
                  onClick={() => setSelectedMemory(memory)}
                  className={`flex flex-col text-left bg-bg-surface border rounded-xl p-4 transition-all hover:shadow-md group ${selectedMemory?.id === memory.id ? 'border-primary ring-1 ring-primary' : 'border-border-base'}`}
                >
                  <div className="flex items-center justify-between mb-3">
                    <div className="flex items-center space-x-2">
                      <div className={`w-2 h-2 rounded-full ${getConfidenceColor(memory.confidence)}`}></div>
                      <span className="text-[10px] font-bold uppercase tracking-tighter text-text-muted">v{memory.version}</span>
                    </div>
                    <span className="text-[10px] text-text-muted font-mono">{memory.id.substring(0, 8)}</span>
                  </div>
                  <p className="text-sm text-text-base line-clamp-3 mb-4 flex-1">
                    {memory.content}
                  </p>
                  <div className="w-full bg-bg-base h-1 rounded-full overflow-hidden">
                    <div 
                      className={`h-full ${getConfidenceColor(memory.confidence)} transition-all duration-500`}
                      style={{ width: `${memory.confidence * 100}%` }}
                    ></div>
                  </div>
                </button>
              ))}
            </div>
          )}
        </div>
      </div>

      {/* Details Sidebar */}
      {selectedMemory && (
        <div className="w-96 border-l border-border-base bg-bg-surface flex flex-col animate-in slide-in-from-right duration-300">
          <div className="p-6 border-b border-border-base flex justify-between items-center flex-none">
            <h2 className="font-bold text-text-base text-lg">Memory Details</h2>
            <button onClick={() => setSelectedMemory(null)} className="text-text-muted hover:text-text-base p-1">
              <svg className="w-5 h-5" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" /></svg>
            </button>
          </div>
          
          <div className="flex-1 overflow-y-auto p-6 space-y-8">
            <section className="space-y-3">
              <h3 className="text-xs font-semibold text-text-muted uppercase tracking-wider">Content</h3>
              <div className="bg-bg-base p-4 rounded-xl text-sm leading-relaxed text-text-base border border-border-base shadow-inner">
                {selectedMemory.content}
              </div>
            </section>

            <section className="space-y-4">
              <div className="flex justify-between items-center">
                <label htmlFor="confidence-slider" className="text-xs font-semibold text-text-muted uppercase tracking-wider">Confidence Evolution</label>
                <span className={`text-sm font-bold ${getConfidenceColor(selectedMemory.confidence).replace('bg-', 'text-')}`}>
                  {(selectedMemory.confidence * 100).toFixed(1)}%
                </span>
              </div>
              <input 
                id="confidence-slider"
                name="confidence-slider"
                type="range"
                min="0"
                max="1"
                step="0.01"
                value={selectedMemory.confidence}
                onChange={(e) => handleFeedback(selectedMemory.id, parseFloat(e.target.value))}
                className="w-full h-2 bg-bg-base rounded-lg appearance-none cursor-pointer accent-primary"
              />
              <p className="text-[10px] text-text-muted italic">Manual adjustment will increment memory version and update semantic search weighting.</p>
            </section>

            <section className="space-y-3">
              <h3 className="text-xs font-semibold text-text-muted uppercase tracking-wider">Provenance</h3>
              <div className="bg-bg-base p-4 rounded-xl border border-border-base font-mono text-[10px] space-y-2 overflow-x-auto">
                {Object.entries(selectedMemory.provenance).map(([key, val]) => (
                  <div key={key} className="flex justify-between">
                    <span className="text-primary">{key}:</span>
                    <span className="text-text-base">{JSON.stringify(val)}</span>
                  </div>
                ))}
              </div>
            </section>
          </div>
        </div>
      )}
    </div>
  );
};

export default MemoryInspector;
