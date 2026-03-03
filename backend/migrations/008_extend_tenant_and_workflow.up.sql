-- Add branding to tenants
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS logo_svg TEXT;
ALTER TABLE tenants ADD COLUMN IF NOT EXISTS brand_title TEXT;

-- Add hierarchy to workflows
ALTER TABLE workflows ADD COLUMN IF NOT EXISTS parent_id UUID REFERENCES workflows(id);
ALTER TABLE workflows ADD COLUMN IF NOT EXISTS element_type TEXT NOT NULL DEFAULT 'workflow' CHECK (element_type IN ('workflow', 'element', 'detail'));

-- Index for tree traversal
CREATE INDEX IF NOT EXISTS idx_workflows_parent_id ON workflows(parent_id);
