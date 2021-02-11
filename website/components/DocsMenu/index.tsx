import React, { useState, useEffect, useContext, createContext } from 'react';
import Link from 'next/link';
import { Collapse, List, ListItem, ListItemText, makeStyles, Theme } from '@material-ui/core';
import Scrollspy from 'react-scrollspy';
import { useRouter } from 'next/router';
import ExpandMore from '@material-ui/icons/ExpandMore';
import ExpandLess from '@material-ui/icons/ExpandLess';

const specLinks = {
  Abstract: 'abstract',
  Terminology: 'terminology',
  'Preload Header': 'preload-header',
  'Using Preload Link Relations': 'using-preload-link-relations',
  'Fields Header': 'fields-header',
  Selectors: 'selectors',
  'Extended JSON Pointer': 'extended-json-pointer',
  'Query Parameters': 'query-parameters',
  'Computing Links Server-Side': 'computing-links-server-side',
  'Security Considerations': 'security-considerations',
  'IANA Considerations': 'iana-considerations',
  'Implementation Status': 'implementation-status',
  'Vulcain Gateway Server': 'vulcain-gateway-server',
  'Helix Vulcain Filters': 'helix-vulcain-filters',
};

interface DocsMenuContextType {
  selectedLink?: string;
}

const DocsMenuContext = createContext<DocsMenuContextType>({});

const useStyles = makeStyles<Theme>((theme) => ({
  root: {
    width: '100%',
  },
  nested: {
    paddingLeft: theme.spacing(4),
  },
}));

interface MenuLinkProps {
  href: string;
  text: string;
  className?: string;
  scroll?: boolean;
}

const MenuLink: React.ComponentType<MenuLinkProps> = ({ href, text, className, scroll = true }) => {
  const { selectedLink } = useContext(DocsMenuContext);

  return (
    <Link href={href} passHref scroll={scroll}>
      <ListItem selected={href === selectedLink} button component="a" className={className}>
        <ListItemText primary={text} />
      </ListItem>
    </Link>
  );
};

const DocsMenu: React.ComponentType = () => {
  const router = useRouter();
  const { asPath } = router;
  const classes = useStyles();
  const [openSpec, setOpenSpec] = useState(asPath.includes('/spec'));
  const [openGateway, setOpenGateway] = useState(asPath.includes('/gateway'));
  const [selectedLink, setSelectedLink] = useState(asPath);

  useEffect(() => {
    setSelectedLink(asPath);
  }, [asPath]);

  return (
    <DocsMenuContext.Provider value={{ selectedLink }}>
      <List component="nav" aria-labelledby="nested-list-subheader" className={classes.root}>
        <MenuLink text="Getting started" href="/docs" />

        <ListItem button onClick={() => setOpenSpec(!openSpec)}>
          <ListItemText primary="Specification" />
          {openSpec ? <ExpandLess /> : <ExpandMore />}
        </ListItem>
        <Collapse in={openSpec} timeout="auto" mountOnEnter unmountOnExit>
          <Scrollspy
            items={Object.values(specLinks)}
            onUpdate={(elem) => {
              if (elem) {
                window.history.replaceState({}, '', `#${elem.id}`);
                setSelectedLink(`/spec/vulcain#${elem.id}`);
              }
            }}
            componentTag={List}
            component="div"
            disablePadding
          >
            {Object.entries(specLinks).map(([k, v]) => (
              <MenuLink text={k} href={`/spec/vulcain#${v}`} scroll={false} key={v} className={classes.nested} />
            ))}
          </Scrollspy>
        </Collapse>

        <ListItem button onClick={() => setOpenGateway(!openGateway)}>
          <ListItemText primary="Gateway" />
          {openGateway ? <ExpandLess /> : <ExpandMore />}
        </ListItem>
        <Collapse in={openGateway} timeout="auto" unmountOnExit>
          <List component="div" disablePadding>
            <MenuLink text="Install" href="/docs/gateway/install" className={classes.nested} />
            <MenuLink text="Configuration" href="/docs/gateway/config" className={classes.nested} />
            <MenuLink text="OpenAPI" href="/docs/gateway/openapi" className={classes.nested} />
          </List>
        </Collapse>

        <MenuLink text="Comparison with GraphQL and Other API Formats" href="/docs/graphql" />

        <MenuLink
          text="Using GraphQL as Query Language for Vulcain"
          href="/docs/graphql#using-graphql-as-query-language-for-vulcain"
        />

        <MenuLink text="Vulcain for Caddy" href="/docs/caddy" />

        <MenuLink text="Cache Considerations" href="/docs/cache" />

        <MenuLink text="Help" href="/docs/help" />

        <MenuLink text="Prior Art" href="/docs/prior-art" />
      </List>
    </DocsMenuContext.Provider>
  );
};

export default DocsMenu;
