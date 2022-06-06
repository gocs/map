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
    let graph = await d3.json('land.json', { cache: "no-store" });
    this.nodes = graph.nodes;
    this.circle = graph.circle;
    this.options = graph.options;
    this.line = graph.line;
  }
};

(async _ => {
  if (!window["WebSocket"]) return;

  var conn = new WebSocket("wss://" + document.location.host + "/updaterws");
  conn.onopen = async (_) => {

    await map.newMap();

    let map_d3 = d3.select('svg g');
    let zoom = d3.zoom()
      .scaleExtent([1, 20])
      .translateExtent([[0, 0], [width, height]])
      .on('zoom', (e, _) => map_d3.attr('transform', e.transform));

    let voronoi = d3.Delaunay.from(map.nodes, d => d.x, d => d.y)
      .voronoi([0, 0, width, height]);

    map_d3.node();
    map_d3.call(zoom).on("dblclick.zoom", null);

    let drag = d3.drag()
      .on("start", (_, d) => circle.filter(p => p === d)
        .raise()
        .attr("stroke", "black")
        .attr("stroke-width", map.circle.stroke_width))
      .on("drag", (event, d) => (d.x = event.x, d.y = event.y))
      .on("end", (_, d) => {
        conn.send(JSON.stringify({ action: "set", ...d }));
        return circle.filter(p => p === d).attr("stroke", "none")
      })
      .on("start.update drag.update end.update", update);

    let hover4circumcircleE = {
      mouseover: (_, d) => map_d3.selectAll("circle").filter(p => p === d).raise().attr("r", 2 * map.circle.radius),
      mouseout: (_, d) => map_d3.selectAll("circle").filter(p => p === d).raise().attr("r", map.circle.radius)
    };

    // render
    let mesh = map_d3
      .append("path")
      .attr("fill", "black")
      .attr("stroke", map.line.color)
      .attr("stroke-width", map.line.stroke_width)
      .attr("d", voronoi.render());

    let cell = map_d3
      .append("g")
      .attr("fill", "none")
      .attr("pointer-events", "all")
      .selectAll("path")
      .data(map.nodes, d => d.id).join("path")
      .on("mouseover", hover4circumcircleE.mouseover)
      .on("mouseout", hover4circumcircleE.mouseout)
      .on("dblclick", (e, d) => conn.send(JSON.stringify({ action: "add", x: e.clientX, y: e.clientY, type: d.type })))
      .attr("fill", (d, _) => map.options.colors[d.type])
      .attr("d", (_, i) => voronoi.renderCell(i));

    let circle = map_d3
      .selectAll('circle')
      .data(map.nodes, d => d.id)
      .enter()
      .append('circle')
      .attr('r', map.circle.radius)
      .attr('data-id', d => d.id)
      .attr('cx', d => d.x)
      .attr('cy', d => d.y)
      .attr("fill", "black")
      .on("dblclick", (_, d) => conn.send(JSON.stringify({ action: "del", id: d.id })))
      .on("mouseover", hover4circumcircleE.mouseover)
      .on("mouseout", hover4circumcircleE.mouseout)
      .call(drag);

    function update() {
      voronoi = d3.Delaunay.from(map.nodes, d => d.x, d => d.y).voronoi([0, 0, width, height]);
      cell
        .attr("fill", (d, _) => map.options.colors[d.type])
        .attr("d", (_, i) => voronoi.renderCell(i));
      mesh.attr("d", voronoi.render());
      circle.attr("cx", d => d.x).attr("cy", d => d.y).attr("fill", "black");
    }

    conn.onclose = evt => console.log("onclose evt:", evt);
    conn.onerror = err => console.error("onerror err:", err);
    conn.onmessage = async evt => {
      let data = JSON.parse(evt.data);

      if (!data) return; // no business
      await map.newMap();

      switch (data.action) {
        case "add":
          map_d3
            .selectAll('circle')
            .data(map.nodes, d => `${d.id}`)
            .enter()
            .append('circle')
            .attr('r', map.circle.radius)
            .attr('data-id', d => `${d.id}`)
            .attr('cx', d => d.x)
            .attr('cy', d => d.y);
          break;
        case "del":
          map_d3
            .selectAll('circle')
            .data(map.nodes, d => `${d.id}`)
            .exit()
            .remove();
          break;
        case "set":
          map_d3
            .selectAll('circle')
            .data(map.nodes, d => `${d.id}`)
            .attr('data-id', d => `${d.id}`)
            .attr("cx", d => d.x).attr("cy", d => d.y);
          break;
      }

      voronoi = d3.Delaunay.from(map.nodes, d => d.x, d => d.y).voronoi([0, 0, width, height]);
      cell
        .attr("fill", (d, _) => map.options.colors[d.type])
        .attr("d", (_, i) => voronoi.renderCell(i));
      mesh.attr("d", voronoi.render());
      circle.attr("cx", d => d.x).attr("cy", d => d.y);
    };
  };
})();
