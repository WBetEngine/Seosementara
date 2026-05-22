(function (w, d) {
  var script = d.currentScript;
  var siteKey = (script && script.getAttribute("data-site")) || "";
  var collectURL = (script && script.getAttribute("data-collect")) || "/collect";

  function uuid() {
    if (w.crypto && crypto.randomUUID) return crypto.randomUUID();
    return "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx".replace(/[xy]/g, function (c) {
      var r = (Math.random() * 16) | 0;
      return (c === "x" ? r : (r & 0x3) | 0x8).toString(16);
    });
  }

  function readCookie(name) {
    var m = d.cookie.match(
      new RegExp("(?:^|; )" + name.replace(/([.$?*|{}()[\]\\/+^])/g, "\\$1") + "=([^;]*)")
    );
    return m ? decodeURIComponent(m[1]) : "";
  }

  function setCookie(name, value, maxAge) {
    d.cookie =
      name +
      "=" +
      encodeURIComponent(value) +
      "; path=/; max-age=" +
      maxAge +
      "; SameSite=Lax";
  }

  function ensureFbp() {
    var v = readCookie("_fbp");
    if (v) return v;
    v = "fb.1." + Math.floor(Date.now() / 1000) + "." + Math.floor(Math.random() * 1e10);
    setCookie("_fbp", v, 7776000);
    return v;
  }

  function fbclidFromURL() {
    try {
      return new URL(w.location.href).searchParams.get("fbclid") || "";
    } catch (e) {
      return "";
    }
  }

  function ensureFbc() {
    var v = readCookie("_fbc");
    if (v) return v;
    var clid = fbclidFromURL();
    if (!clid) return "";
    v = "fb.1." + Math.floor(Date.now() / 1000) + "." + clid;
    setCookie("_fbc", v, 7776000);
    return v;
  }

  function send(eventName, props) {
    props = props || {};
    if (!props.event_id) props.event_id = uuid();
    var payload = {
      event: eventName,
      event_id: props.event_id,
      url: w.location.href,
      site_key: siteKey,
      fbp: ensureFbp(),
      fbc: ensureFbc(),
      fbclid: fbclidFromURL(),
      props: props
    };
    if (props.email) payload.email = props.email;
    if (props.phone) payload.phone = props.phone;
    if (props.first_name) payload.first_name = props.first_name;
    if (props.last_name) payload.last_name = props.last_name;
    if (props.external_id) payload.external_id = props.external_id;
    if (props.country) payload.country = props.country;
    if (props.value != null) payload.props.value = props.value;
    if (props.currency) payload.props.currency = props.currency;
    if (props.order_id) payload.props.order_id = props.order_id;
    if (props.content_ids) payload.props.content_ids = props.content_ids;

    var body = JSON.stringify(payload);
    if (navigator.sendBeacon) {
      navigator.sendBeacon(collectURL, new Blob([body], { type: "application/json" }));
    } else {
      fetch(collectURL, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: body,
        keepalive: true
      });
    }
    return props.event_id;
  }

  w.sseo = w.sseo || {};
  w.sseo.track = send;
  w.sseo.trackLead = function (opts) {
    return send("lead", opts || {});
  };
  w.sseo.trackPurchase = function (opts) {
    opts = opts || {};
    if (!opts.value || !opts.currency) {
      console.warn("[sseo] Purchase butuh value dan currency (Meta Plan/25)");
    }
    return send("purchase", opts);
  };

  send("page_view", { path: w.location.pathname });
})(window, document);
