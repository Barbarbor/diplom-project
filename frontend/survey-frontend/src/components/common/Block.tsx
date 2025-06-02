import React from 'react';

interface BlockProps {
  children: React.ReactNode;
  className?: string;
}

export const Block = ({ children, className }: BlockProps) => {
  return (
    <div className={`mb-6 bg-white shadow-md rounded-lg p-4 ${className}`}>
   
      {children}
    </div>
  );
};