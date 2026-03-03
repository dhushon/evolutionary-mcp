-- Grounding Rules (formerly Anchors)
-- High-priority context pointers for AI stability
CREATE TABLE IF NOT EXISTS grounding_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    tenant_id UUID NOT NULL REFERENCES tenants(id),
    workflow_id UUID REFERENCES workflows(id), -- Optional link to specific workflow element
    name TEXT NOT NULL,
    content TEXT NOT NULL,
    embedding VECTOR(384), -- Semantic vector for grounding lookup
    is_global BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_grounding_rules_tenant ON grounding_rules(tenant_id);
CREATE INDEX IF NOT EXISTS idx_grounding_rules_workflow ON grounding_rules(workflow_id);
-- Vector index for fast retrieval
CREATE INDEX IF NOT EXISTS idx_grounding_rules_embedding ON grounding_rules USING ivfflat (embedding vector_l2_ops) WITH (lists = 100);
