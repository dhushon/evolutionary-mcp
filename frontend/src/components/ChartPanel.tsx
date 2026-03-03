import React from 'react';
import { BarChart, Bar, XAxis, YAxis, CartesianGrid, Tooltip, Legend, ResponsiveContainer } from 'recharts';

interface ChartPanelProps {
    data: any[];
}

const ChartPanel: React.FC<ChartPanelProps> = ({ data }) => {
    const processedData = data.reduce((acc, item) => {
        const userId = item.userId;
        if (acc[userId]) {
            acc[userId].count += 1;
        } else {
            acc[userId] = { userId, count: 1 };
        }
        return acc;
    }, {});

    const chartData = Object.values(processedData);

    return (
        <ResponsiveContainer width="100%" height={300}>
            <BarChart data={chartData}>
                <CartesianGrid strokeDasharray="3 3" />
                <XAxis dataKey="userId" />
                <YAxis />
                <Tooltip />
                <Legend />
                <Bar dataKey="count" fill="#8884d8" />
            </BarChart>
        </ResponsiveContainer>
    );
};

export default ChartPanel;
