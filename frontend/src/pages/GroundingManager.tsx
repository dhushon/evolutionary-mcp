import React, { useState } from 'react';
import { 
  useGroundingRules, 
  useCreateGroundingRule, 
  useUpdateGroundingRule, 
  useDeleteGroundingRule 
} from '../hooks/useGrounding';
import { GroundingRule } from '../types';

const GroundingManager: React.FC = () => {
  const { data: rules, isLoading } = useGroundingRules();
  const createMutation = useCreateGroundingRule();
  const updateMutation = useUpdateGroundingRule();
  const deleteMutation = useDeleteGroundingRule();

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingRule, setEditingRule] = useState<Partial<GroundingRule> | null>(null);

  const handleSave = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!editingRule) return;

    try {
      if (editingRule.id) {
        await updateMutation.mutateAsync({ id: editingRule.id, rule: editingRule });
      } else {
        await createMutation.mutateAsync(editingRule);
      }
      setIsModalOpen(false);
      setEditingRule(null);
    } catch (err) {
      console.error("Failed to save rule:", err);
    }
  };

  return (
    <div className="p-6 space-y-6">
      <div className="flex justify-between items-center">
        <div>
          <h1 className="text-2xl font-bold text-text-base">Grounding Rules</h1>
          <p className="text-sm text-text-muted">Foundational truths and reasoning stability vectors.</p>
        </div>
        <button 
          onClick={() => { setEditingRule({ name: '', content: '', is_global: false }); setIsModalOpen(true); }}
          className="bg-primary text-white px-4 py-2 rounded-lg hover:bg-primary-hover transition-all"
        >
          + New Rule
        </button>
      </div>

      {isLoading ? (
        <div className="flex justify-center py-12">
          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary"></div>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {rules?.map(rule => (
            <div key={rule.id} className="bg-bg-surface border border-border-base rounded-xl p-5 shadow-sm hover:shadow-md transition-all group">
              <div className="flex justify-between items-start mb-3">
                <h3 className="font-bold text-text-base truncate pr-4">{rule.name}</h3>
                <div className="flex space-x-1 opacity-0 group-hover:opacity-100 transition-opacity">
                  <button 
                    onClick={() => { setEditingRule(rule); setIsModalOpen(true); }}
                    className="p-1 hover:text-primary transition-colors"
                  >
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" /></svg>
                  </button>
                  <button 
                    onClick={() => { if(confirm("Delete rule?")) deleteMutation.mutate(rule.id); }}
                    className="p-1 hover:text-red-500 transition-colors"
                  >
                    <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor"><path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" /></svg>
                  </button>
                </div>
              </div>
              <p className="text-sm text-text-muted line-clamp-3 mb-4 leading-relaxed">
                {rule.content}
              </p>
              <div className="flex items-center justify-between mt-auto pt-4 border-t border-border-base/50">
                <span className={`text-[10px] uppercase tracking-wider font-bold px-2 py-0.5 rounded-full ${rule.is_global ? 'bg-purple-100 text-purple-700 dark:bg-purple-900/30 dark:text-purple-400' : 'bg-blue-100 text-blue-700 dark:bg-blue-900/30 dark:text-blue-400'}`}>
                  {rule.is_global ? 'Global' : 'Tenant'}
                </span>
                <span className="text-[10px] text-text-muted">
                  {new Date(rule.updated_at).toLocaleDateString()}
                </span>
              </div>
            </div>
          ))}
        </div>
      )}

      {isModalOpen && (
        <div className="fixed inset-0 z-50 flex items-center justify-center p-4 bg-black/50 backdrop-blur-sm animate-in fade-in duration-200">
          <div className="bg-bg-surface w-full max-w-lg rounded-2xl shadow-xl border border-border-base animate-in zoom-in-95 duration-200">
            <form onSubmit={handleSave} className="p-6 space-y-4">
              <h2 className="text-xl font-bold text-text-base">
                {editingRule?.id ? 'Edit Grounding Rule' : 'New Grounding Rule'}
              </h2>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-text-base mb-1">Rule Name</label>
                  <input 
                    type="text"
                    required
                    value={editingRule?.name || ''}
                    onChange={e => setEditingRule(prev => ({ ...prev!, name: e.target.value }))}
                    className="w-full bg-bg-base border border-border-base rounded-lg px-4 py-2 outline-none focus:ring-2 focus:ring-primary transition-all"
                    placeholder="e.g. Identity Verification Protocol"
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-text-base mb-1">Content / Instruction</label>
                  <textarea 
                    required
                    rows={6}
                    value={editingRule?.content || ''}
                    onChange={e => setEditingRule(prev => ({ ...prev!, content: e.target.value }))}
                    className="w-full bg-bg-base border border-border-base rounded-lg px-4 py-2 outline-none focus:ring-2 focus:ring-primary transition-all resize-none"
                    placeholder="Define the factual grounding or reasoning constraint..."
                  />
                </div>
                <div className="flex items-center space-x-2">
                  <input 
                    type="checkbox"
                    id="is_global"
                    checked={editingRule?.is_global || false}
                    onChange={e => setEditingRule(prev => ({ ...prev!, is_global: e.target.checked }))}
                    className="rounded border-border-base text-primary focus:ring-primary"
                  />
                  <label htmlFor="is_global" className="text-sm font-medium text-text-base">Global Rule (system-wide)</label>
                </div>
              </div>
              <div className="flex justify-end space-x-3 pt-4 border-t border-border-base">
                <button 
                  type="button"
                  onClick={() => setIsModalOpen(false)}
                  className="px-4 py-2 text-sm font-medium text-text-muted hover:text-text-base transition-colors"
                >
                  Cancel
                </button>
                <button 
                  type="submit"
                  disabled={createMutation.isPending || updateMutation.isPending}
                  className="bg-primary text-white px-6 py-2 rounded-lg font-medium hover:bg-primary-hover shadow-sm transition-all disabled:opacity-50"
                >
                  Save Rule
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default GroundingManager;
