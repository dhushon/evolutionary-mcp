import React from 'react';

interface DataCardProps {
  title: string;
  subtitle?: string;
  content: string;
  footer?: React.ReactNode;
  confidence?: number;
}

export const DataCard: React.FC<DataCardProps> = ({ title, subtitle, content, footer, confidence }) => {
  return (
    <div className="bg-white dark:bg-slate-800 rounded-xl border border-slate-200 dark:border-slate-700 shadow-sm hover:shadow-md transition-shadow overflow-hidden">
      <div className="p-5">
        <div className="flex justify-between items-start mb-2">
          <h3 className="font-semibold text-slate-900 dark:text-white truncate">{title}</h3>
          {confidence !== undefined && (
            <span className={`text-xs font-medium px-2 py-1 rounded-full ${confidence > 0.8 ? 'bg-green-100 text-green-700' : 'bg-amber-100 text-amber-700'}`}>
              {Math.round(confidence * 100)}% match
            </span>
          )}
        </div>
        {subtitle && <p className="text-xs text-slate-500 dark:text-slate-400 mb-3">{subtitle}</p>}
        <p className="text-sm text-slate-600 dark:text-slate-300 line-clamp-3">{content}</p>
      </div>
      {footer && <div className="px-5 py-3 bg-slate-50 dark:bg-slate-800/50 border-t border-slate-100 dark:border-slate-700">{footer}</div>}
    </div>
  );
};