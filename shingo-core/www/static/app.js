// SSE connection for live updates
(function() {
  let es;

  function connect() {
    es = new EventSource('/events');

    es.addEventListener('order-update', function(e) {
      // Refresh order-related content
      const el = document.querySelector('[data-sse="orders"]');
      if (el) htmx.trigger(el, 'sse-refresh');
      // Also refresh dashboard stats
      const dash = document.querySelector('[data-sse="dashboard"]');
      if (dash) htmx.trigger(dash, 'sse-refresh');
    });

    es.addEventListener('inventory-update', function(e) {
      const el = document.querySelector('[data-sse="nodestate"]');
      if (el) htmx.trigger(el, 'sse-refresh');
    });

    es.addEventListener('node-update', function(e) {
      const el = document.querySelector('[data-sse="nodes"]');
      if (el) htmx.trigger(el, 'sse-refresh');
    });

    es.addEventListener('system-status', function(e) {
      const data = JSON.parse(e.data);
      if (data.fleet !== undefined) {
        const el = document.getElementById('fleet-status');
        if (el) {
          el.className = 'health ' + (data.fleet === 'connected' ? 'health-ok' : 'health-fail');
        }
      }
      if (data.messaging !== undefined) {
        const el = document.getElementById('msg-status');
        if (el) {
          el.className = 'health ' + (data.messaging === 'connected' ? 'health-ok' : 'health-fail');
        }
      }
      if (data.redis !== undefined) {
        const el = document.getElementById('redis-status');
        if (el) {
          el.className = 'health ' + (data.redis === 'connected' ? 'health-ok' : 'health-fail');
        }
      }
    });

    es.addEventListener('robot-update', function(e) {
      var robots = JSON.parse(e.data);
      var grid = document.getElementById('robot-grid');
      if (!grid) return;

      var seen = {};
      robots.forEach(function(r) {
        seen[r.vehicle_id] = true;
        var tile = grid.querySelector('[data-name="' + r.vehicle_id + '"]');
        if (!tile) {
          // Create new tile
          tile = document.createElement('div');
          tile.className = 'robot-tile robot-' + r.state;
          tile.setAttribute('onclick', 'openRobotModal(this)');
          tile.innerHTML =
            '<div class="robot-name">' + r.vehicle_id +
            (r.charging ? '<span class="robot-charging" title="Charging">&#9889;</span>' : '') +
            '</div>' +
            '<div class="robot-battery" title="Battery: ' + r.battery + '%">' +
            '<div class="robot-battery-fill" style="width:' + r.battery + '%"></div>' +
            '</div>';
          grid.appendChild(tile);
        } else {
          // Update tile class
          tile.className = 'robot-tile robot-' + r.state;
          // Update battery bar
          var fill = tile.querySelector('.robot-battery-fill');
          if (fill) fill.style.width = r.battery + '%';
          var batDiv = tile.querySelector('.robot-battery');
          if (batDiv) batDiv.title = 'Battery: ' + r.battery + '%';
          // Update charging indicator
          var nameDiv = tile.querySelector('.robot-name');
          if (nameDiv) {
            var chgSpan = nameDiv.querySelector('.robot-charging');
            if (r.charging && !chgSpan) {
              chgSpan = document.createElement('span');
              chgSpan.className = 'robot-charging';
              chgSpan.title = 'Charging';
              chgSpan.innerHTML = '&#9889;';
              nameDiv.appendChild(chgSpan);
            } else if (!r.charging && chgSpan) {
              chgSpan.remove();
            }
          }
        }
        // Update data attributes
        tile.dataset.name = r.vehicle_id;
        tile.dataset.state = r.state;
        tile.dataset.ip = r.ip || '';
        tile.dataset.model = r.model || '';
        tile.dataset.map = r.map || '';
        tile.dataset.battery = r.battery;
        tile.dataset.charging = r.charging;
        tile.dataset.station = r.station || '';
        tile.dataset.lastStation = r.last_station || '';
        tile.dataset.available = r.available;
        tile.dataset.connected = r.connected;
        tile.dataset.blocked = r.blocked;
        tile.dataset.emergency = r.emergency;
        tile.dataset.processing = r.processing;
        tile.dataset.error = r.error;
        tile.dataset.x = r.x.toFixed(1);
        tile.dataset.y = r.y.toFixed(1);
        tile.dataset.angle = r.angle.toFixed(1);

        // Update modal if open for this robot
        if (typeof currentRobotVehicle !== 'undefined' && currentRobotVehicle === r.vehicle_id) {
          var modal = document.getElementById('robot-modal');
          if (modal && modal.classList.contains('active')) {
            openRobotModal(tile);
          }
        }
      });

      // Remove stale tiles
      var tiles = grid.querySelectorAll('.robot-tile');
      tiles.forEach(function(tile) {
        if (!seen[tile.dataset.name]) {
          tile.remove();
        }
      });

      // Update robot count
      var countEl = document.getElementById('robot-count');
      if (countEl) {
        countEl.textContent = robots.length + ' robots';
      }

      // Show/hide empty state
      var emptyCard = grid.nextElementSibling;
      if (robots.length === 0 && !grid.children.length) {
        grid.style.display = 'none';
      } else {
        grid.style.display = '';
      }

      // Reapply filter
      if (typeof filterRobots === 'function') {
        filterRobots();
      }
    });

    es.addEventListener('debug-log', function(e) {
      if (typeof window.debugAppendRow === 'function') {
        var entry = JSON.parse(e.data);
        window.debugAppendRow(entry);
      }
    });

    es.onerror = function() {
      es.close();
      setTimeout(connect, 3000);
    };
  }

  connect();
})();
