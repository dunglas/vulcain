import React from 'react';
import { Divider, Drawer, Box, Theme, makeStyles } from '@material-ui/core';
import DocsMenu from '../DocsMenu';

const useStyles = makeStyles<Theme>((theme) => ({
  root: {
    width: '80%',
  },
  header: {
    backgroundColor: theme.palette.primary.main,
    color: '#fff',
    display: 'flex',
    flexDirection: 'column',
    alignItems: 'center',
    padding: theme.spacing(2),
  },
  accountIcon: {
    fontSize: '8rem',
  },
  menu: {
    width: '100%',
  },
}));

interface MobileMenuProps {
  open?: boolean;
  onClose?: () => void;
}

const MobileMenu: React.ComponentType<MobileMenuProps> = ({ open, onClose }) => {
  const classes = useStyles();
  return (
    <Drawer open={open} anchor="right" onClose={onClose} PaperProps={{ classes: { root: classes.root } }}>
      <Divider />
      <Box flex={1} alignItems="flex-start" display="flex">
        <DocsMenu />
      </Box>
    </Drawer>
  );
};

export default MobileMenu;
