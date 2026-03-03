export type WorkflowStatus = 'draft' | 'active' | 'archived';
export type ElementType = 'workflow' | 'element' | 'detail';

export interface Workflow {
  id: string;
  workflow_id: string;
  tenant_id: string;
  name: string;
  description: string;
  status: WorkflowStatus;
  version: number;
  is_latest: boolean;
  parent_id?: string;
  element_type: ElementType;
  input_schema: Record<string, any>;
  output_schema?: Record<string, any>;
  created_by: string;
  created_at: string;
  updated_at: string;
}

export interface WorkflowUpdatePayload extends Partial<Workflow> {
  save_as_new_version?: boolean;
}

export interface Tenant {
  id: string;
  name: string;
  domain: string;
  logo_svg?: string;
  brand_title?: string;
  created_at: string;
  updated_at: string;
}

export interface GroundingRule {
  id: string;
  tenant_id: string;
  workflow_id?: string;
  name: string;
  content: string;
  is_global: boolean;
  created_at: string;
  updated_at: string;
}

export interface HealthStatus {
  service: string;
  status: string;
  timestamp: string;
  version: string;
}
