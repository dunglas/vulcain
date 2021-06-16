import React from 'react';
import { METHODS } from '../data/methods';
import CustomTabs from './CustomTabs';

interface MethodSelectorProps {
  method?: string;
  onMethodChange: (val: string) => void;
}

const MethodSelector: React.ComponentType<MethodSelectorProps> = ({ onMethodChange, method }) => {
  const tabs = Object.keys(METHODS).map((key) => ({ value: key, label: METHODS[key].label }));

  return <CustomTabs tabs={tabs} value={method} onChange={onMethodChange} />;
};

export default MethodSelector;
