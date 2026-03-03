export interface Workflow {
    id: string;
    workflow_id: string;
    tenant_id: string;
    name: string;
    description: string;
    status: 'draft' | 'active' | 'archived';
    version: number;
    is_latest: boolean;
    input_schema: Record<string, any>;
    output_schema?: Record<string, any>;
    created_by: string;
    created_at: string;
    updated_at: string;
}