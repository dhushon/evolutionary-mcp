import { Handle, Position } from 'reactflow';

const TransformNode = ({ data }: { data: any }) => {
  return (
    <div style={{ border: '1px solid #777', padding: 10, borderRadius: 5, background: 'white' }}>
      <div>Transform Node</div>
      <input type="text" defaultValue={data.transform || ''} placeholder="Enter transform" style={{ marginTop: 5 }} />
      <Handle type="target" position={Position.Top} />
      <Handle type="source" position={Position.Bottom} />
    </div>
  );
};

export default TransformNode;
