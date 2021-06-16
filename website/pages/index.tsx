import React from 'react';
import Page from '../components/Page';
import Main from '../components/home/Main';
import Support from '../components/home/Support';
import Features from '../components/home/Features';
import References from '../components/home/References';

const Index: React.ComponentType = () => (
  <Page>
    <Main />
    <Features />
    <References />
    <Support />
  </Page>
);

export default Index;
