import React, { useState, useEffect } from 'react';
import { AgGridReact } from 'ag-grid-react';
import { CellValueChangedEvent } from 'ag-grid-community';
import 'ag-grid-community/styles/ag-grid.css';
import 'ag-grid-community/styles/ag-theme-alpine.css';

interface DataTableProps {
    rowData: any[];
    onCellValueChanged: (event: CellValueChangedEvent) => void;
}

const DataTable: React.FC<DataTableProps> = ({ rowData, onCellValueChanged }) => {
    const [columnDefs, setColumnDefs] = useState<any[]>([]);

    useEffect(() => {
        if (rowData && rowData.length > 0) {
            const keys = Object.keys(rowData[0]);
            const defs = keys.map(key => ({
                field: key,
                sortable: true,
                filter: true,
                editable: true,
            }));
            setColumnDefs(defs);
        }
    }, [rowData]);

    return (
        <div className="ag-theme-alpine" style={{ height: 400, width: '100%' }}>
            <AgGridReact
                rowData={rowData}
                columnDefs={columnDefs}
                onCellValueChanged={onCellValueChanged}
            />
        </div>
    );
};

export default DataTable;
