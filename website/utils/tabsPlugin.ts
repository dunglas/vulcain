function tabsPlugin() {
  function isTabBlockCode(node) {
    return node?.children?.length === 1 && node.children[0].value === '[tabs]';
  }

  function isEndTabBlockCode(node) {
    return node?.children?.length === 1 && node.children[0].value === '[/tabs]';
  }

  function transformer(tree, file) {
    let tabBlockCode;
    let depth;
    const TAB_TITLE_REGEX = /---(.*)---((.|\n|\r)*)/i;

    for (let child of tree.children) {
      // start of tab block code
      if (isTabBlockCode(child)) {
        tabBlockCode = child;
        tabBlockCode.children = [];
        tabBlockCode.tabs = [];
        tabBlockCode.type = 'tabs';
        // end of tab block code
      } else if (isEndTabBlockCode(child)) {
        tabBlockCode = null;
        depth = null;
        child.toDelete = true;
      } else if (tabBlockCode) {
        // if child is heading node: check for new tab
        if (child.type === 'heading') {
          if (!depth) depth = child.depth;
          // new tab
          if (child.depth === depth) {
            const tabTitle = child.children[0]?.value;
            tabBlockCode.tabs.push(tabTitle);
            tabBlockCode.children.push({
              type: 'tab',
              children: [],
              value: tabTitle,
            });
            // if no tab: put child on last tab node
          } else {
            tabBlockCode.children.length &&
              tabBlockCode.children[tabBlockCode.children.length - 1].children.push({ ...child });
          }
          child.toDelete = true;
          // else: put child on last tab node
        } else {
          tabBlockCode.children.length &&
            tabBlockCode.children[tabBlockCode.children.length - 1].children.push({ ...child });
          child.toDelete = true;
        }
      }
    }
    // remove duplicated and useless nodes
    tree.children = tree.children.filter((node) => !node.toDelete);
  }
  return transformer;
}

export default tabsPlugin;
