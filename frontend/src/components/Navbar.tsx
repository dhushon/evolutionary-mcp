import React from 'react';
import { useTenant } from '../hooks/useWorkflows';

const Navbar: React.FC = () => {
  const { data: tenant } = useTenant();

  return (
    <nav className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700 p-4 shadow-sm z-20">
      <div className="container mx-auto flex justify-between items-center">
        <div className="flex items-center space-x-3">
          <h1 className="text-xl font-bold text-gray-900 dark:text-white">
            {tenant?.brand_title || 'Data Management & Workflow'}
          </h1>
        </div>
        <div className="flex-1 max-w-md mx-8">
          <div className="relative">
            <span className="absolute inset-y-0 left-0 pl-3 flex items-center text-gray-400">
              <svg className="w-4 h-4" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
              </svg>
            </span>
            <input 
              type="search" 
              placeholder="Search workflows, memories..." 
              className="w-full bg-gray-100 dark:bg-gray-700 border-none rounded-lg pl-10 pr-4 py-2 text-sm focus:ring-2 focus:ring-primary transition-all" 
            />
          </div>
        </div>
        <div className="flex items-center space-x-4">
          <div className="text-sm font-medium text-gray-700 dark:text-gray-300">
            {tenant?.name || 'User Profile'}
          </div>
          <div className="w-8 h-8 rounded-full bg-primary/20 flex items-center justify-center text-primary font-bold text-xs">
            {tenant?.name?.charAt(0) || 'U'}
          </div>
        </div>
      </div>
    </nav>
  );
};

export default Navbar;
