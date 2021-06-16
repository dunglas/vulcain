import React, { useState } from 'react';
import TabContext from '@material-ui/lab/TabContext';
import { Box } from '@material-ui/core';
import CustomTabs from '../CustomTabs';
import Card from '@material-ui/core/Card';

interface TabsListProps {
  tabs: string[];
}

const TabsList: React.ComponentType<TabsListProps> = ({ tabs, children }) => {
  const [value, setValue] = useState(tabs[0]);

  const handleChange = (newValue) => {
    setValue(newValue);
  };

  return (
    <TabContext value={value}>
      <Box pb={2}>
        <Card square elevation={5}>
          <Box py={2}>
            <CustomTabs value={value} tabs={tabs} onChange={handleChange} />
          </Box>
          <Box textAlign="center">{children}</Box>
        </Card>
      </Box>
    </TabContext>
  );
};

export default TabsList;
