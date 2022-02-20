let width = window.innerWidth, height = window.innerHeight;
let points = 10;

let map = {
  nodes: [],
  circle: {
    polygon: 0,
    radius: 0,
    stroke_width: 0,
  },
  options: { colors: {} },
  line: {},
  async newMap() {
    let graph = await d3.json('land.json');
    this.nodes = graph.nodes;
    this.circle = graph.circle;
    this.options = graph.options;
    this.line = graph.line;
  },
  add(mesh) {
    this.nodes.push(mesh);
  },
  set(mesh) {
    this.nodes[mesh.id] = mesh;
  },
  del(id) {
    // this.nodes = this.nodes.filter((d) => d.id != id);
    this.nodes.splice(id, 1);
  },
};

(async _ => {
  await map.newMap();

  let map_d3 = d3
    .select('svg g');

  map_d3
    .selectAll('circle')
    .data(map.nodes, d => d.id)
    .enter()
    .append('circle')
    .attr('r', map.circle.radius)
    .style('fill', d => map.options.colors[d.type])
    .attr('data-id', d => d.id)
    .attr('cx', d => d.x)
    .attr('cy', d => d.y);

  // window.addEventListener("click", evt => {
  //   console.log("click target:", evt.target);
  // });

  window.addEventListener("dblclick", evt => {
    if (!evt.target.__data__) {
      map.add({ id: 1, x: evt.clientX, y: evt.clientY, type: "untyped" });
      

      map_d3
        .selectAll('circle')
        .data(map.nodes, d => d.id)
        .enter()
        .append('circle')
        .attr('r', map.circle.radius)
        .style('fill', d => map.options.colors[d.type])
        .attr('data-id', d => d.id)
        .attr('cx', d => d.x)
        .attr('cy', d => d.y);
      return;
    }
    map.del(evt.target.__data__.id);

    map_d3
      .selectAll('circle')
      .data(map.nodes, d => d.id)
      .exit()
      .remove();
  });

  // window.addEventListener("mouseover", evt => {
  //   console.log("mouseover target:", evt.target);
  // })
  // window.addEventListener("mouseout", evt => {
  //   console.log("mouseout target:", evt.target);
  // })

})();


// (async _ => {
//   if (!window["WebSocket"]) return;

//   var conn = new WebSocket("ws://" + document.location.host + "/updaterws");
//   conn.onclose = function (evt) {
//     console.log("onclose evt:", evt);
//   };
//   conn.onmessage = async function (evt) {
//     console.log("onmessage evt:", JSON.parse(evt.data));
//     voronoi = d3.Delaunay.from(circles, d => d.x, d => d.y).voronoi([0, 0, width, height]);
//     cell
//       .attr("fill", (d, _) => graph.options.colors[d.type])
//       .attr("d", (_, i) => voronoi.renderCell(i));
//     mesh.attr("d", voronoi.render());
//     circle.attr("cx", d => d.x).attr("cy", d => d.y);
//   };

//   let graph = await d3.json('land.json');
//   // data
//   let circles = graph.nodes;
//   let map_group = d3.select('svg g');
//   let voronoi = d3.Delaunay.from(circles, d => d.x, d => d.y)
//     .voronoi([0, 0, width, height]);

//   // events
//   // hovering a cell will show its circumcircle
//   let hover4circumcircleE = {
//     mouseover: (_, d) => map_group.selectAll("circle").filter(p => p === d).raise().attr("r", graph.circle.radius),
//     mouseout: (_, d) => map_group.selectAll("circle").filter(p => p === d).raise().attr("r", 0)
//   };

//   // double clicking a cell will create a new cell
//   let dblclick2addnodeE = {
//     dblclick: (e, d) => {
//       let latestID = circles.reduce((p, c) => p > c.id ? p : c.id, 0);
//       let data = { id: latestID, x: e.clientX, y: e.clientY, type: d.type };
//       conn.send(JSON.stringify({ action: "addcell", ...data }));

//       circles.push(data);
//       // update();
//       // fetch('/appendcell', {
//       //   method: 'POST',
//       //   headers: { 'Content-Type': 'application/json' },
//       //   body: JSON.stringify({ x: e.clientX, y: e.clientY, type: d.type }),
//       // }).then(_ => _/**window.location.reload() */).catch(err => console.error('Error dblclick2addnodeE:', err))
//     }
//   };

//   let dblclick2deletenodeE = {
//     dblclick: (e, d) => {
//       conn.send(JSON.stringify({ action: "delcell", id: d.id }));
//       fetch('/deletecell', {
//         method: 'DELETE',
//         headers: { 'Content-Type': 'application/json' },
//         body: JSON.stringify({ id: d.id }),
//       }).then(_ => _/**window.location.reload() */).catch(err => console.error('Error dblclick2deletenodeE:', err))
//     }
//   };

//   let drag = d3.drag()
//     .on("start", (_, d) => circle.filter(p => p === d)
//       .raise()
//       .attr("stroke", "black")
//       .attr("stroke-width", graph.circle.stroke_width))
//     .on("drag", (event, d) => (d.x = event.x, d.y = event.y))
//     .on("end", (_, d) => {
//       conn.send(JSON.stringify({ action: "movecell", ...d }));
//       fetch('/landjson', {
//         method: 'POST',
//         headers: { 'Content-Type': 'application/json' },
//         body: JSON.stringify(d),
//       }).catch(err => console.error('Error:', err));
//       return circle.filter(p => p === d).attr("stroke", "none")
//     })
//     .on("start.update drag.update end.update", update);

//   let zoom = d3.zoom()
//     .scaleExtent([1, 5])
//     .translateExtent([[0, 0], [width, height]])
//     .on('zoom', (e, _) => d3.select('svg g').attr('transform', e.transform));

//   // render
//   let mesh = map_group
//     .append("path")
//     .attr("fill", "black")
//     .attr("stroke", graph.line.color)
//     .attr("stroke-width", graph.line.stroke_width)
//     .attr("d", voronoi.render());

//   let cell = map_group
//     .append("g")
//     .attr("fill", "none")
//     .attr("pointer-events", "all")
//     .selectAll("path")
//     .data(circles).join("path")
//     .on("mouseover", hover4circumcircleE.mouseover)
//     .on("mouseout", hover4circumcircleE.mouseout)
//     .on("dblclick", dblclick2addnodeE.dblclick)
//     .attr("fill", (d, i) => graph.options.colors[d.type])
//     .attr("d", (d, i) => voronoi.renderCell(i));

//   let circle = map_group
//     .selectAll("circle")
//     .data(circles).join("circle")
//     .attr("cx", d => d.x)
//     .attr("cy", d => d.y)
//     .attr("r", 0)
//     .on("mouseover", hover4circumcircleE.mouseover)
//     .on("mouseout", hover4circumcircleE.mouseout)
//     .on("dblclick", dblclick2deletenodeE.dblclick)
//     .call(drag)
//     .attr("fill", "black");

//   function update() {
//     voronoi = d3.Delaunay.from(circles, d => d.x, d => d.y).voronoi([0, 0, width, height]);
//     cell
//       .attr("fill", (d, _) => graph.options.colors[d.type])
//       .attr("d", (_, i) => voronoi.renderCell(i));
//     mesh.attr("d", voronoi.render());
//     circle.attr("cx", d => d.x).attr("cy", d => d.y);
//   }

//   map_group.node();
//   map_group.call(zoom).on("dblclick.zoom", null);
// })()