package hotreload

const clientScript = `
<script>
(function() {
    if (window.__dadeHotReload) return;
    window.__dadeHotReload = true;

    console.log('[Dade] Initializing hot reload (polling)...');

    var pollInterval = 1000;
    var lastHash = '';

    function checkForUpdate() {
        fetch('/css/output.css?' + Date.now())
            .then(function(response) {
                return response.text();
            })
            .then(function(css) {
                var hash = hashCode(css);
                if (lastHash && hash !== lastHash) {
                    console.log('[Dade] CSS changed, reloading page...');
                    window.location.reload();
                }
                lastHash = hash;
            })
            .catch(function(error) {
                console.error('[Dade] Polling error:', error);
            });
    }

    function hashCode(str) {
        var hash = 0;
        for (var i = 0; i < str.length; i++) {
            var char = str.charCodeAt(i);
            hash = ((hash << 5) - hash) + char;
            hash = hash & hash;
        }
        return hash.toString();
    }

    // Start polling
    console.log('[Dade] Starting to poll for changes every', pollInterval, 'ms');
    checkForUpdate();
    setInterval(checkForUpdate, pollInterval);
})();
</script>
`

// GetClientScript returns the client-side SSE script as a string
func GetClientScript() string {
	return clientScript
}
