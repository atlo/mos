import { radioController } from './radio';
import { switchController } from './switch';
import { sliders } from './slider';

/* eslint-disable */
const getColor = type => {
  switch (type) {
    case 'sajtó':
      return '#ff5155';
    case 'internet':
      return '#009dc5';
    case 'televízió':
      return '#7ca929';
    case 'rádió':
      return '#ffc455';
    case 'érdekeltség':
      return '#fff';
    default:
      return '#aaa';
  }
};

const getLinkPoints = (el) => {
  const a = el.source.x || 0;
  const b = el.source.y || 0;
  const c = el.target.x || 0;
  const d = el.target.y || 0;

  if ((el.source.dashed || el.target.dashed)) {
    return `${a},${b} ${c},${d}`;  
  }

  let r = el.source.size;
  let R = el.target.size;

  if (!(el.source.large || el.target.large)) { r = 2; R = 2 }
  if (r - R === 0) { r -= 0.1 }

  const xp = (c * r - a * R) / (r - R);
  const yp = (d * r - b * R) / (r - R);

  const xt1 = (r * r * (xp - a) + r * (yp - b) * Math.sqrt((xp - a) * (xp - a) + (yp - b) * (yp - b) - r * r)) / ((xp - a) * (xp - a) + (yp - b) * (yp - b)) + a;
  const yt1 = (r * r * (yp - b) - r * (xp - a) * Math.sqrt((xp - a) * (xp - a) + (yp - b) * (yp - b) - r * r)) / ((xp - a) * (xp - a) + (yp - b) * (yp - b)) + b;

  const xt2 = (r * r * (xp - a) - r * (yp - b) * Math.sqrt((xp - a) * (xp - a) + (yp - b) * (yp - b) - r * r)) / ((xp - a) * (xp - a) + (yp - b) * (yp - b)) + a;
  const yt2 = (r * r * (yp - b) + r * (xp - a) * Math.sqrt((xp - a) * (xp - a) + (yp - b) * (yp - b) - r * r)) / ((xp - a) * (xp - a) + (yp - b) * (yp - b)) + b;

  const xt3 = (R * R * (xp - c) + R * (yp - d) * Math.sqrt((xp - c) * (xp - c) + (yp - d) * (yp - d) - R * R)) / ((xp - c) * (xp - c) + (yp - d) * (yp - d)) + c;
  const yt3 = (R * R * (yp - d) - R * (xp - c) * Math.sqrt((xp - c) * (xp - c) + (yp - d) * (yp - d) - R * R)) / ((xp - c) * (xp - c) + (yp - d) * (yp - d)) + d;

  const xt4 = (R * R * (xp - c) - R * (yp - d) * Math.sqrt((xp - c) * (xp - c) + (yp - d) * (yp - d) - R * R)) / ((xp - c) * (xp - c) + (yp - d) * (yp - d)) + c;
  const yt4 = (R * R * (yp - d) + R * (xp - c) * Math.sqrt((xp - c) * (xp - c) + (yp - d) * (yp - d) - R * R)) / ((xp - c) * (xp - c) + (yp - d) * (yp - d)) + d;

  return `${a},${b} ${xt1},${yt1} ${xt3},${yt3} ${c},${d} ${xt4},${yt4} ${xt2},${yt2}`;
}

const updateLinkPoints = (d) => { 
  const result = getLinkPoints(d);
  if (result.split('NaN').length > 1) {
    return '0,0 0,0 0,0 0,0 0,0 0,0';
  }
  return result;
}

class GraphDetail {
  constructor(data) {
    this.clicked = data.clicked;
    this.hash = data.hash;
    this.links = data.links;
    
    let ownInd = 0, medInd = 0, opInd = 0, intrestInd = 0;
    this.group = data.group.map((id) => {
      let type = null;

      if (this.hash[id].hasOwnProperty('hun_non_hun')) {
        type = 'owner';
        return Object.assign({ objType: type, i: ownInd++ }, this.hash[id]);
      };
      if (this.hash[id].hasOwnProperty('news_non_news')) {
        type = 'media';
        return Object.assign({ objType: type, i: medInd++ }, this.hash[id]);
      };
      if (this.hash[id].hasOwnProperty('data')) {
        type = 'operator';
        return Object.assign({ objType: type, i: opInd++ }, this.hash[id]);
      };
      if (this.hash[id].hasOwnProperty('Name')) {
        type = 'interest';
        return Object.assign({ objType: type, i: intrestInd++ }, this.hash[id]);
      };
    });

    this.detailedGraphContainer = document.querySelector('[data-overlay-graph]');
    this.detailedGraphContainer.innerHTML = '';
    this.tooltip = d3.select('body')
      .append('div')
      .attr('class', 'tooltip')
      .style('opacity', 0);

    const dimensions = this.detailedGraphContainer.getBoundingClientRect();
    this.w = dimensions.width;
    this.h = dimensions.height;

    this.svg = d3.select('[data-overlay-graph]')
      .append('svg')
      .attr('width', this.w)
      .attr('height', this.h);

    this.link = this.svg.append('g')
      .attr('class', 'links')
      .selectAll('line')
      .data(this.links)
      .enter().append('svg:polyline')
      .attr('class', 'link')
      .attr('points',  (d) => {
        const link = Object.assign(d, {});
        link.source.x =  this.getColumn(this.group.find(item => item['_id'] === d.source['_id']));
        link.source.y = this.getRow(this.group.find(item => item['_id'] === d.source['_id']));
        link.target.x =  this.getColumn(this.group.find(item => item['_id'] === d.target['_id']));
        link.target.y = this.getRow(this.group.find(item => item['_id'] === d.target['_id']));

        const result = getLinkPoints(link);
        if (result.split('NaN').length > 1) {
          return '0,0 0,0 0,0 0,0 0,0 0,0';
        }
        return result;
      })
      .attr('fill', '#ddd')
      .attr('stroke', d => {
        return (
          d.source.dashed || d.target.dashed
        ) ? '#fff' : 'none'
      })
      .attr('stroke-dasharray', d => {
        return (
          d.source.dashed || d.target.dashed
        ) ? '2' : '0'
      })
      .attr('opacity', d => {
        return (
          d.source.dashed || d.target.dashed
        ) ? 1 : .1
      });
    

    this.nodes = this.svg.append('g')
      .attr('class', 'nodes')
      .selectAll('circle')
      .data(this.group)
      .enter()
      .append('circle')
      .attr('cx', this.getColumn.bind(this))
      .attr('cy', this.getRow.bind(this))
      .attr('r', (d) => (d.size || 5))
      .attr('fill', (d) => { return getColor(d.type || 'default') })
      .attr('stroke-width', (d) => {
        return d['_id'] === this.clicked ? "3" : "none"
      })
      .attr('stroke', (d) => {
        return d['_id'] === this.clicked ? "#fff" : "none"
      })
      .on('mouseover', this.mouseover.bind(this))
      .on('mouseout', this.mouseout.bind(this));

    this.svg
      .attr('height', this.nodes._parents[0].getBoundingClientRect().height + 80);

    document.querySelector('[data-overlay-title]').innerText =
      this.hash[this.clicked].name || 'N/A';
  }

  mouseover(d) {
    const hovered = this.group.find(item => item['_id'] === d['_id']);
    const x =  this.getColumn(hovered);
    const y = this.getRow(hovered);
    const size = hovered.size;
    const offs = this.svg._groups[0][0].getBoundingClientRect();
    
    this.tooltip
      .style('opacity', 1);

    this.tooltip
      .html(d.name)
      .attr('style', `left: ${x + offs.left}px; top: ${y + offs.top - size - 10}px;`);
  };

  mouseout() {
    this.tooltip
      .style('opacity', 0)
  };

  getColumn(d) {
    if (d.objType === 'media') {
      return (this.w - 20) / 8 * 7;
    }
    if (d.objType === 'operator') {
      return (this.w - 20) / 8 * 5;
    }
    if (d.objType === 'owner') {
      return (this.w - 20) / 8 * 3;
    }
    if (d.objType === 'interest') {
      return (this.w - 20) / 8 * 1;
    }
  }

  getRow(d) {
    return (40 + d['i'] * 50);
  }
}

class G {
  constructor(selector, w, h, data, max_profits, startYear) {
    this.w = w;
    this.h = h;
    this.dataAll = Object.assign({}, data);
    this.selector = selector;
    this.currentYear = startYear;
    this.currentDetailYear = startYear;
    this.max_profits = max_profits;
    
    switchController.subscribe((value, checked) => {
      this.onSwitchChange(value, checked);
    })

    radioController.subscribe((value) => {
      this.onRadioChange(value);
    });

    sliders.forEach(slider => {
      if (slider.id === 'main') {
        slider.subscribe(this.onMainSliderChange.bind(this));
      } else {
        slider.subscribe(this.onDetailSliderChange.bind(this));
      }
    });

    this.init();
  }

  init() {
    if ((this.data || {}).nodes) this.data.nodes = [];
    if ((this.data || {}).links) this.data.links = [];
    this.data = JSON.parse(JSON.stringify(this.dataAll[this.currentYear]));
    this.groups = this.getGroups();
    this.clickedNodeId = null;
    this.nodeHash = this.getNodeHash();

    this.svg = d3.select(this.selector)
      .html('')
      .append('svg')
      .attr('width', this.w)
      .attr('height', this.h);

    const padding = 50;
    const groupPositions = [];
    this.groups = Object.entries(this.groups).reduce((acc, [year, group]) => {
      acc[year] = group.sort((a, b) => { return b.length - a.length })
      return acc;
    }, {});

    const valid = (x, y, size) => {
      return groupPositions.every(circle => {
        const sum = circle.size + size;
        const dist = Math.sqrt((circle.x - x) * (circle.x - x) + (circle.y - y) * (circle.y - y));
        return sum < dist;
      });
    };

    this.groups[this.currentYear].forEach((g) => {
      const size = Math.min(g.length * 15, this.h / 15);
      let x = 0;
      let y = 0;
      let maxIter = 500;
      do {
        x = Math.random() * (this.w - padding * 2) + padding;
        y = Math.random() * (this.h - padding * 2) + padding;
        maxIter--;
      } while (!valid(x, y, size) && maxIter > 0);

      groupPositions.push({ x, y, size });
    });

    this.simulation = d3.forceSimulation()
      .force('link', d3.forceLink().id((d) => (d['_id'])))
      .force('charge', d3.forceManyBody().strength(-1))
      .force('collision', d3.forceCollide().radius(d => {
        if (d.hasOwnProperty('hun_non_hun')) {
          return (d.size * 3 || 10);
        };
        return (d.size * 2 || 10);
      }))
      .force('center', d3.forceCenter(this.w / 2, this.h / 2))
      .force('x', d3.forceX().x(d => {
        const found = this.groups[this.currentYear].findIndex(group => {
          return group.some(linkedId => linkedId === d['_id']);
        });
        return (groupPositions[found] || { x: this.w / 2 }).x;
      }))
      .force('y', d3.forceY().y(d => {
        const found = this.groups[this.currentYear].findIndex(group => {
          return group.some(linkedId => linkedId === d['_id']);
        });
        return (groupPositions[found] || { y: this.h / 2 }).y;
      }));

    this.link = this.svg.append('g')
      .attr('class', 'links')
      .selectAll('line')
      .data(this.data.links)
      .enter().append('svg:polyline')
      .attr('class', 'link')
      .attr('points', updateLinkPoints)
      .attr('fill', '#ddd')
      .attr('stroke', d => {
        return (
          this.nodeHash[this.currentDetailYear][d.source].dashed || 
          this.nodeHash[this.currentDetailYear][d.target].dashed
        ) ? '#fff' : 'none'
      })
      .attr('stroke-dasharray', d => {
        return (
          this.nodeHash[this.currentDetailYear][d.source].dashed || 
          this.nodeHash[this.currentDetailYear][d.target].dashed
        ) ? '2' : '0'
      })
      .attr('opacity', d => {
        return (
          this.nodeHash[this.currentDetailYear][d.source].dashed || 
          this.nodeHash[this.currentDetailYear][d.target].dashed
        ) ? 1 : .1
      });

    this.node = this.svg.append('g')
      .attr('class', 'nodes')
      .selectAll('circle')
      .data(this.data.nodes)
      .enter()
      .append('circle')
      .attr('r', (d) => (d.size || 5))
      .attr('fill', (d) => { return getColor(d.type || 'default') })
      .call(d3.drag()
        .on('start', this.dragstarted.bind(this))
        .on('drag', this.dragged.bind(this))
        .on('end', this.dragended.bind(this))
      )
      .on('mouseover', this.mouseover.bind(this))
      .on('mouseout', this.mouseout.bind(this))
      .on('click', this.click.bind(this));

    this.tooltip = d3.select('body')
      .append('div')
      .attr('class', 'tooltip')
      .style('opacity', 0);

    this.simulation
      .nodes(this.data.nodes)
      .on('tick', () => this.ticked(this.w, this.h));

    this.simulation.force('link')
      .links(this.data.links);
  }

  getGroups() {
    const groups = {};
    Object.entries(Object.assign({}, this.dataAll)).forEach(([year, data]) => {
      data.links.forEach(link => {
        if (!groups[year]) {
          groups[year] = [];
        }

        let firstGroup = null;
        let secondGroup = null;
        for(let i = 0; i < groups[year].length; i++) {
          if ((groups[year][i] || []).some(id => ((id || null) === link.source ))) {
            firstGroup = i;
          }
          if ((groups[year][i] || []).some(id => ((id || null) === link.target ))) {
            secondGroup = i;
          }
        }
        if (firstGroup !== null && secondGroup !== null) {
          if (firstGroup !== secondGroup) {
            groups[year][firstGroup] = [
              ...groups[year][firstGroup],
              ...groups[year][secondGroup],
            ];
            groups[year].splice(secondGroup, 1);
          }
        } else if (firstGroup !== null && secondGroup === null) {
          groups[year][firstGroup] = [
            ...groups[year][firstGroup],
            link.target,
          ];
  
        } else if (firstGroup === null && secondGroup !== null) {
          groups[year][secondGroup] = [
            ...groups[year][secondGroup],
            link.source,
          ];
  
        } else {
          groups[year].push([
            link.source,
            link.target,
          ]);
  
        }
      });
    });
    return groups;
  }

  getNodeHash() {
    let nodeHash= {};

    Object.entries(this.dataAll).forEach(([year, data]) => {
      if (!nodeHash[year]) {
        nodeHash[year] = {};
      }
      data.nodes.forEach((node) => {
        nodeHash[year][node['_id']] = node;
      });
    });
    return nodeHash;
  }

  update(data) {
    this.data = data;
  }

  ticked(w, h) {
    this.link
      .attr('points', (d) => { 
        const result = getLinkPoints(d);
        if (result.split('NaN').length > 1) {
          return '0,0 0,0 0,0 0,0 0,0 0,0';
        }
        return result;
      });

    this.node
      .attr('cx', function(d) { return d.x = Math.max(d.size || 10, Math.min(w - d.size || 10, d.x)); })
      .attr('cy', function(d) { return d.y = Math.max(d.size || 10, Math.min(h - d.size || 10, d.y)); });
  }

  dragstarted(d) {
    if (!d3.event.active) this.simulation.alphaTarget(0.3).restart();
    d.fx = d.x;
    d.fy = d.y;
  }

  dragged(d) {
    d.fx = d3.event.x;
    d.fy = d3.event.y;
  }

  dragended(d) {
    if (!d3.event.active) this.simulation.alphaTarget(0);
    d.fx = null;
    d.fy = null;
  }

  mouseover(d) {
    this.tooltip
      .style('opacity', 1);
    
    let html = d.name;
    if (d.hasOwnProperty('data')) {
      html = `<div class="tooltip__inner">
        <div class="tooltip__content">${d.name}</div>
        <div class="tooltip__label">Cím</div>
        <div class="tooltip__content">${d[this.currentYear].address !== '' ? d[this.currentYear].address : 'N/A'}</div>
        <div class="tooltip__label">Nettó árbevétel</div>
        <div class="tooltip__content">${d[this.currentYear].netto_profit / 1000 + 'M Forint'}</div>
        <div class="tooltip__label">Adózott eredmény</div>
        <div class="tooltip__content">${d[this.currentYear].taxed_profit / 1000 + 'M Forint'}</div>
        <div class="tooltip__label">Üzemi eredmény</div>
        <div class="tooltip__content">${d[this.currentYear].operating_profit / 1000 + 'M Forint'}</div>
      </div>`
    }

    this.tooltip
      .html(html)
      .attr('class', d.y > this.h / 2 ? 'tooltip tooltip--down' : 'tooltip tooltip--up')
      .attr('style',
        `left: ${d.x + 312}px; top: ${d.y > this.h / 2 ? 
          d.y - d.size - 10 : 
          d.y + d.size + 10 + this.tooltip._groups[0][0].getBoundingClientRect().height}px;`
      )
  };

  mouseout() {
    this.tooltip
      .style('opacity', 0)
  };

  click(d) {
    this.clickedNodeId = d['_id'];
    this.currentDetailYear = this.currentYear;
    this.onGraphDetailChange();
  }

  onGraphDetailChange() {
    document.querySelector('[data-overlay="graph"]').classList.add('overlay--open');
    sliders.find(slider => {
      return slider.id !== 'main';
    }).setYear(this.currentDetailYear);
    
    const currentGroup = this.groups[this.currentDetailYear].find(
      group => group.some(item => item === this.clickedNodeId)
    );

    if (!currentGroup) {
      document.querySelector('[data-overlay-graph]').classList.add('overlay__graph--no-data');
      return;
    } else {
      document.querySelector('[data-overlay-graph]').classList.remove('overlay__graph--no-data');
    }

    const currentLinks = this.dataAll[this.currentDetailYear].links.filter(
      link => {
        return (
          currentGroup.includes(link.target['_id'] || link.target
          ) || 
          currentGroup.includes(link.source['_id'] || link.source
          )
        )
      }
    ).map(item => ({
      source: Object.assign({}, typeof item.source === 'number' ? 
        this.nodeHash[this.currentDetailYear][item.source] : item.source),
      target: Object.assign({}, typeof item.target === 'number' ? 
        this.nodeHash[this.currentDetailYear][item.target] : item.target),
    }));

    this.detailedGraph = new GraphDetail({
      clicked: this.clickedNodeId,
      group: currentGroup,
      links: currentLinks,
      hash: this.nodeHash[this.currentDetailYear],
    });
  }

  onSwitchChange(type, checked) {
    this.node
      .attr('fill', (d) => {
        if (!checked) {
          return getColor(d.type || 'default');
        }

        let newColor = '#aaa';
        switch (type) {
          case 'hu-nonhu':
            if (d.hasOwnProperty('hun_non_hun')) {
              newColor = d.hun_non_hun ? '#ff5155': '#009dc5';
            }
            return newColor;
            break;
          case 'left-right':
            if (d.hasOwnProperty('left_right')) {
              newColor = d.left_right ? '#ff5155': '#009dc5';
            }
            return newColor;
            break;
          case 'media-nonmedia':
            if (d.hasOwnProperty('news_non_news')) {
              newColor = d.news_non_news ? '#ff5155': '#009dc5';
            }
            return newColor;
            break;
          default:
            return getColor(d.type || 'default');
            break;
        }
      });
  }

  onRadioChange(type) {
    this.node
      .attr('r', (d) => {
        if (d && d.hasOwnProperty('data')) {  
          let newSize = d.size;
          switch (type) {
            case 'net':
              newSize = Math.max(d[this.currentYear].netto_profit / 
                this.max_profits[this.currentYear].netto_profit, 0) * 20 + 5;
              d.size = newSize < 0 ? 5 : newSize;
              return newSize < 0 ? 5 : newSize;
              break;
            case 'tax':
              newSize = Math.max(d[this.currentYear].taxed_profit / 
                this.max_profits[this.currentYear].taxed_profit, 0) * 20 + 5;
              d.size = newSize < 0 ? 5 : newSize;
              return newSize < 0 ? 5 : newSize;
              break;
            case 'operational':
              newSize = Math.max(d[this.currentYear].operating_profit / 
                this.max_profits[this.currentYear].operating_profit, 0) * 20 + 5;
              d.size = newSize < 0 ? 5 : newSize;
              return newSize < 0 ? 5 : newSize;
              break;
            default:
              break;
          }
        } else {
          return d.size;
        };
      });
    
    this.link
      .attr('points', updateLinkPoints);
  }

  onMainSliderChange(newYear) {
    if (newYear === this.currentYear) return;
    this.currentYear = newYear;
    this.currentDetailYear = newYear;

    this.init();
  }

  onDetailSliderChange(newDetailYear) {
    if (newDetailYear === this.currentDetailYear) return;
    this.currentDetailYear = newDetailYear;

    this.onGraphDetailChange();
  }
}



export default class Graph {
  constructor() {
    const startYear = 1998;
    // Change api path for local development to use localhost:4000
    // Otherwise use relative path
    // Example: http://localhost:4000/media
    Promise.all([
      fetch('/media').then(res => res.json()),
      fetch('/connections').then(res => res.json()),
      fetch('/operators').then(res => res.json()),
      fetch('/owners').then(res => res.json()),
      fetch('/interests').then(res => res.json()),
    ]).then(res => {
      const media = res[0];
      let connections = res[1];
      let operators = res[2];
      const owners = res[3];
      const interests = res[4];
      const allYears = Object.entries(operators[0].data).map(([key]) => (parseInt(key, 10)));
      let max_profits = {};

      connections = connections.map((conn) => (Object.assign(conn, conn.data)));
      operators = operators.map((operator) => (Object.assign(operator, operator.data)));

      allYears.forEach((year) => {
        max_profits[year] = operators.reduce((result, operator) => {
          result.netto_profit = ((operator[year] || {}).netto_profit || 0) > result.netto_profit ? 
          ((operator[year] || {}).netto_profit || 0) : result.netto_profit;
          result.taxed_profit = ((operator[year] || {}).taxed_profit || 0) > result.taxed_profit ? 
          ((operator[year] || {}).taxed_profit || 0) : result.taxed_profit;
          result.operating_profit = ((operator[year] || {}).operating_profit || 0) > result.operating_profit ? 
          ((operator[year] || {}).operating_profit || 0) : result.operating_profit;
          return result;
        }, { netto_profit: 0, taxed_profit: 0, operating_profit: 0 })
      });

      const data = allYears.reduce((data, value) => {
        data[value] = {
          nodes: [
            ...owners
              .filter(owner => {
                return connections.some(connection => {
                  return (connection[value] || {}).ownerIds
                    && connection[value].ownerIds.some(val => parseInt(val) === owner['_id'])
                })
              })
              .map(owner => (Object.assign(owner, {
                size: 5,
              }))),
            ...media
              .filter(media => {
                return connections.some(connection => {
                  return (connection[value] || {}).mediaIds
                    && connection[value].mediaIds.some(val => parseInt(val) === media['_id'])
                })
              })
              .map(media => (Object.assign(media, {
                size: 5,
                large: true,
              }))),
            ...operators
              .filter(operator => {
                return connections.some(connection => {
                  return connection[value] && connection['_id'] === operator['_id']
                })
              })
              .map(operator => {
                if (Math.max(operator[value].netto_profit / max_profits[value].netto_profit, 0) >= 1 && value === startYear) {
                  console.log(operator, value, max_profits[value], Math.max(operator[value].netto_profit / max_profits[value].netto_profit, 0) * 20 + 5);
                }
                return Object.assign(operator, {
                  size: Math.max(operator[value].netto_profit / max_profits[value].netto_profit, 0) * 20 + 5
                })
              }),
            ...interests
              .filter(interest => {
                return operators.some(operator => {
                  return ((operator[value] || {}).interests || []).some(intr => intr === interest.Id)
                })
              })
              .map(interest => (Object.assign(interest, {
                size: 5,
                dashed: true,
                type: 'érdekeltség',
                name: interest.Name,
                _id: interest.Id,
              })))
          ],
          links: []
        };
        return data;
      }, {});

      console.log(max_profits[startYear], operators);

      allYears.forEach((value) => {
        operators
          .filter(operator => {
            return connections.some(connection => {
              return connection[value] && connection['_id'] === operator['_id']
            })
          })
          .forEach(operator => {
            if (operator[value]) {
              const interestIds = operator[value].interests;
              interestIds.forEach(interestId => {
                data[value].links.push({
                  source: operator['_id'],
                  target: interestId,
                });
              });
            }
          });

        connections.forEach(connection => {
          if (connection[value]) {
            const opId = connection['_id'];
            const medIds = (connection[value].mediaIds || []).filter((item, pos, self) => {
              return self.indexOf(item) == pos;
            });
            const ownIds = (connection[value].ownerIds || []).filter((item, pos, self) => {
              return self.indexOf(item) == pos;
            });

            if (medIds) {
              medIds.forEach(medId => {
                if (data[value].nodes.some(val => val['_id'] === parseInt(medId))) {
                  data[value].links.push({
                    source: opId,
                    target: parseInt(medId),
                  });
                }
              });
            }

            if (ownIds) {
              ownIds.forEach(ownId => {
                if (data[value].nodes.some(val => val['_id'] === parseInt(ownId))) {
                  data[value].links.push({
                    source: opId,
                    target: parseInt(ownId),
                  });
                }
              });
            }
          }
        });
      });

      this.graph = new G(
        '#graph', 
        window.innerWidth - 312, 
        window.innerHeight, 
        Object.assign({}, data), 
        max_profits, 
        startYear // allYears[0]
      );
    }).catch(err => console.log(err));
  }
}
