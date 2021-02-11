import React from 'react';
import { makeStyles, Theme } from '@material-ui/core/styles';
import { Tabs, Tab } from '@material-ui/core';

const useTabsStyles = makeStyles<Theme>((theme) => ({
  root: {
    minHeight: 40,
  },
  flexContainer: {
    display: 'inline-flex',
    position: 'relative',
    zIndex: 1,
    padding: theme.spacing(0.5),
  },
  scroller: {
    marginLeft: 'auto',
    marginRight: 'auto',
    width: 'auto',
    flex: 'unset',
    backgroundColor: theme.palette.type === 'light' ? '#eee' : theme.palette.divider,
    borderRadius: 10,
    padding: 0,
  },
  indicator: {
    top: theme.spacing(0.5),
    bottom: theme.spacing(0.5),
    right: theme.spacing(0.5),
    height: 'auto',
    background: 'none',
    '&:after': {
      content: '""',
      display: 'block',
      position: 'absolute',
      top: 0,
      left: 0,
      right: 0,
      bottom: 0,
      borderRadius: 8,
      backgroundColor: theme.palette.type === 'light' ? '#fff' : theme.palette.action.selected,
      boxShadow: '0 4px 12px 0 rgba(0,0,0,0.16)',
    },
  },
}));
const useTabItemStyles = makeStyles<Theme>((theme) => ({
  root: {
    minHeight: 30,
    minWidth: 100,
  },
  wrapper: {
    textTransform: 'initial',
    fontWeight: theme.typography.fontWeightBold,
  },
}));

type TabValue = {
  value: string;
  label: string;
};

interface CustomTabsProps {
  value?: string;
  onChange: (val: string) => void;
  tabs: (TabValue | string)[];
}

const CustomTabs: React.ComponentType<CustomTabsProps> = ({ onChange, value, tabs }) => {
  const tabsStyles = useTabsStyles();
  const tabItemStyles = useTabItemStyles();

  const formattedTabs = tabs.map((tab) => (typeof tab === 'string' ? { value: tab, label: tab } : tab));

  return (
    <Tabs textColor="primary" centered classes={tabsStyles} value={value} onChange={(e, index) => onChange(index)}>
      {formattedTabs.map((tab) => (
        <Tab classes={tabItemStyles} key={tab.value} disableRipple label={tab.label} value={tab.value} />
      ))}
    </Tabs>
  );
};

export default CustomTabs;
