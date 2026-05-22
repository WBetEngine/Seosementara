(function (w, d) {
  var script = d.currentScript;
  var siteKey = script && script.getAttribute("data-site") || "";
  var collectURL = (script && script.getAttribute("data-collect")) || "/collect";

  function send(eventName, props) {
    var payload = {
      event: eventName,
      event_id: (props && props.event_id) || null,
      url: w.location.href,
      site_key: siteKey,
      fbp: readCookie("_fbp"),
      fbc: readCookie("_fbc"),
      props: props || {}
    };
    var body = JSON.stringify(payload);
    if (navigator.sendBeacon) {
      navigator.sendBeacon(collectURL, body);
    } else {
      fetch(collectURL, { method: "POST", headers: { "Content-Type": "application/json" }, body: body, keepalive: true });
    }
  }

  function readCookie(name) {
    var m = d.cookie.match(new RegExp("(?:^|; )" + name.replace(/([.$?*|{}()[\]\\/+^])/g, "\\$1") + "=([^;]*)"));
    return m ? decodeURIComponent(m[1]) : "";
  }

  w.sseo = w.sseo || {};
  w.sseo.track = send;
  send("page_view", { path: w.location.pathname });
})(window, document);
