let width = window.innerWidth, height = window.innerHeight;

// const zoom = d3.zoom()
//   .scaleExtent([1, 5])
//   .translateExtent([[0, 0], [width, height]])
//   .on('zoom', e => d3.select('svg g')
//     .attr('transform', e.transform));

// d3.json('land.json')
//   .then(graph => {
//     graph.links.forEach(d => {
//       d.source = graph.nodes[d.source];
//       d.target = graph.nodes[d.target];
//     });

//     let points = [];
//     let polygon = graph.edge.polygon;
//     for (let i = 0; i <= polygon; i++) {
//       points.push({ angle: (2 * Math.PI) * i / polygon });
//     }

//     const pointsData = d3.radialArea()
//       .angle(d => d.angle)
//       .outerRadius(_ => graph.edge.radius)(points);

//     d3.select('svg')
//       .call(zoom);
//     let svg = d3.select('svg g');


//     // line
//     let link = svg
//       .attr("class", "link")
//       .selectAll("line")
//       .data(graph.links)
//       .enter().append("line")
//       .style("stroke", d => graph.line[(d.type || "") + "color"])
//       .style("stroke-width", graph.line.width)
//       .attr("x1", d => d.source.x)
//       .attr("y1", d => d.source.y)
//       .attr("x2", d => d.target.x)
//       .attr("y2", d => d.target.y);

//     // edge
//     let node = svg
//       .attr("class", "node")
//       .selectAll("path")
//       .data(graph.nodes)
//       .enter().append('path')
//       .attr('d', pointsData)
//       .attr('fill', graph.edge.color)
//       .attr('transform', d => `translate(${d.x}, ${d.y})`)
//       .call(d3.drag().on("drag", dragged));

//     function dragged(d) {
//       d.subject.x = d.x, d.subject.y = d.y;
//       d3.select(this).attr('transform', `translate(${d.subject.x}, ${d.subject.y})`);
//       link.filter(l => l.source === d.subject).attr("x1", d.subject.x).attr("y1", d.subject.y);
//       link.filter(l => l.target === d.subject).attr("x2", d.subject.x).attr("y2", d.subject.y);
//     }
//   });

const w = window.innerWidth;
const h = (w * 9) / 16;
const canvas = document.createElement("canvas"); 
const context = canvas.getContext("2d"); // DOM.context2d(width, height);

canvas.width = w;
canvas.height = h;

const data = Array(100)
  .fill()
  .map((_, i) => ({ x: (i * w) / 100, y: Math.random() * h }));

const voronoi = d3.Delaunay.from(
  data,
  (d) => d.x,
  (d) => d.y
).voronoi([0, 0, w, h]);

context.clearRect(0, 0, w, h);

context.fillStyle = "black";
context.beginPath();
voronoi.delaunay.renderPoints(context, 1);
context.fill();

context.lineWidth = 1.5;

const segments = voronoi.render().split(/M/).slice(1);
let i = 0;
for (const e of segments) {
  context.beginPath();
  context.strokeStyle = d3.hsl(360 * Math.random(), 0.7, 0.5);
  context.stroke(new Path2D("M" + e));
}
// no yield context.canvas; as we're not on a generator

document.querySelector("#app").appendChild(canvas);