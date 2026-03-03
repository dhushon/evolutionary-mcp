import React, { useState, useEffect } from 'react';
import { Workflow, WorkflowUpdatePayload } from '../types';

interface WorkflowFormProps {
  initialData: Partial<Workflow>;
  onSave: (data: WorkflowUpdatePayload) => void;
  onCancel: () => void;
  isSubmitting?: boolean;
}

const WorkflowForm: React.FC<WorkflowFormProps> = ({ 
  initialData, 
  onSave, 
  onCancel, 
  isSubmitting 
}) => {
  const [formData, setFormData] = useState<Partial<Workflow>>(initialData);
  const [errors, setErrors] = useState<Record<string, string>>({});

  useEffect(() => {
    setFormData(initialData);
  }, [initialData]);

  const validate = () => {
    const newErrors: Record<string, string> = {};
    if (!formData.name) newErrors.name = 'Name is required';
    
    // Validate JSON schemas
    try {
      if (typeof formData.input_schema === 'string') {
        JSON.parse(formData.input_schema);
      }
    } catch (e) {
      newErrors.input_schema = 'Invalid JSON format';
    }

    try {
      if (typeof formData.output_schema === 'string') {
        JSON.parse(formData.output_schema);
      }
    } catch (e) {
      newErrors.output_schema = 'Invalid JSON format';
    }

    setErrors(newErrors);
    return Object.keys(newErrors).length === 0;
  };

  const handleSaveDraft = (e: React.FormEvent) => {
    e.preventDefault();
    if (validate()) {
      onSave({ ...formData, save_as_new_version: false });
    }
  };

  const handlePublishNewVersion = (e: React.FormEvent) => {
    e.preventDefault();
    if (validate()) {
      onSave({ ...formData, save_as_new_version: true, status: 'active' });
    }
  };

  const handleChange = (e: React.ChangeEvent<HTMLInputElement | HTMLTextAreaElement | HTMLSelectElement>) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleJsonChange = (name: 'input_schema' | 'output_schema', value: string) => {
    try {
      const parsed = JSON.parse(value);
      setFormData(prev => ({ ...prev, [name]: parsed }));
      setErrors(prev => ({ ...prev, [name]: '' }));
    } catch (e) {
      // Keep as string for editing, but flag error
      setFormData(prev => ({ ...prev, [name]: value }));
      setErrors(prev => ({ ...prev, [name]: 'Invalid JSON' }));
    }
  };

  return (
    <form className="space-y-6 bg-bg-surface p-6 rounded-xl border border-border-base shadow-sm">
      <div className="space-y-4">
        <div>
          <label className="block text-sm font-medium text-text-base mb-1">Name</label>
          <input
            type="text"
            name="name"
            value={formData.name || ''}
            onChange={handleChange}
            className={`w-full bg-bg-base border ${errors.name ? 'border-red-500' : 'border-border-base'} rounded-lg px-4 py-2 focus:ring-2 focus:ring-primary outline-none transition-all`}
            placeholder="e.g. Data Processor"
          />
          {errors.name && <p className="text-red-500 text-xs mt-1">{errors.name}</p>}
        </div>

        <div>
          <label className="block text-sm font-medium text-text-base mb-1">Description</label>
          <textarea
            name="description"
            value={formData.description || ''}
            onChange={handleChange}
            rows={3}
            className="w-full bg-bg-base border border-border-base rounded-lg px-4 py-2 focus:ring-2 focus:ring-primary outline-none transition-all resize-none"
            placeholder="What does this workflow element do?"
          />
        </div>

        <div className="grid grid-cols-2 gap-4">
          <div>
            <label className="block text-sm font-medium text-text-base mb-1">Type</label>
            <select
              name="element_type"
              value={formData.element_type || 'workflow'}
              onChange={handleChange}
              className="w-full bg-bg-base border border-border-base rounded-lg px-4 py-2 focus:ring-2 focus:ring-primary outline-none transition-all"
            >
              <option value="workflow">Workflow</option>
              <option value="element">Element</option>
              <option value="detail">Detail</option>
            </select>
          </div>
          <div>
            <label className="block text-sm font-medium text-text-base mb-1">Status</label>
            <select
              name="status"
              value={formData.status || 'draft'}
              onChange={handleChange}
              className="w-full bg-bg-base border border-border-base rounded-lg px-4 py-2 focus:ring-2 focus:ring-primary outline-none transition-all"
            >
              <option value="draft">Draft</option>
              <option value="active">Active</option>
              <option value="archived">Archived</option>
            </select>
          </div>
        </div>

        <div>
          <label className="block text-sm font-medium text-text-base mb-1">Input Schema (JSON)</label>
          <textarea
            value={typeof formData.input_schema === 'object' ? JSON.stringify(formData.input_schema, null, 2) : formData.input_schema || '{}'}
            onChange={(e) => handleJsonChange('input_schema', e.target.value)}
            rows={5}
            className={`w-full font-mono text-xs bg-bg-base border ${errors.input_schema ? 'border-red-500' : 'border-border-base'} rounded-lg px-4 py-2 focus:ring-2 focus:ring-primary outline-none transition-all resize-none`}
          />
          {errors.input_schema && <p className="text-red-500 text-xs mt-1">{errors.input_schema}</p>}
        </div>
      </div>

      <div className="flex items-center justify-between pt-4 border-t border-border-base">
        <button
          type="button"
          onClick={onCancel}
          className="text-sm font-medium text-text-muted hover:text-text-base transition-colors"
        >
          Cancel
        </button>
        <div className="flex space-x-3">
          <button
            type="button"
            onClick={handleSaveDraft}
            disabled={isSubmitting}
            className="px-4 py-2 text-sm font-medium text-primary bg-primary/10 rounded-lg hover:bg-primary/20 transition-all disabled:opacity-50"
          >
            Save Draft
          </button>
          <button
            type="button"
            onClick={handlePublishNewVersion}
            disabled={isSubmitting}
            className="px-4 py-2 text-sm font-medium text-white bg-primary rounded-lg hover:bg-primary-hover shadow-sm transition-all disabled:opacity-50"
          >
            Save New Version
          </button>
        </div>
      </div>
    </form>
  );
};

export default WorkflowForm;
