import { Handle, Position } from 'reactflow';

const FilterNode = ({ data }: { data: any }) => {
  return (
    <div style={{ border: '1px solid #777', padding: 10, borderRadius: 5, background: 'white' }}>
      <label htmlFor="filter-condition" style={{ display: 'block', fontSize: '10px', color: '#666' }}>Filter Node</label>
      <input 
        id="filter-condition"
        name="filter-condition"
        type="text" 
        defaultValue={data.condition || ''} 
        placeholder="Enter condition" 
        style={{ marginTop: 5, width: '100%', fontSize: '12px' }} 
      />
      <Handle type="target" position={Position.Top} />
      <Handle type="source" position={Position.Bottom} />
    </div>
  );
};

export default FilterNode;
