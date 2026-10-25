const assert = require('node:assert/strict');
const { readFileSync } = require('node:fs');
const { test } = require('node:test');
const vm = require('node:vm');

function fixture() {
    const listeners = new Map();
    const instances = new Map();
    let dark = false;
    const counts = () => ({
        tagName: 'CANVAS',
        dataset: { chart: 'counts', labels: '["6A"]', datasets: '{"Topped":[2]}' }
    });
    const progression = () => ({
        tagName: 'CANVAS',
        dataset: { chart: 'progression', labels: '["W1","W2","W3"]', progression: '{"6A":[{"Value":3},{"Value":0},{"Value":2}]}' }
    });
    const document = {
        canvases: [counts(), progression()],
        documentElement: { classList: { contains: () => dark } },
        addEventListener: (name, listener) => listeners.set(name, listener),
        querySelectorAll: () => document.canvases
    };
    class Chart {
        static getChart(canvas) { return instances.get(canvas); }
        constructor(canvas, config) {
            assert.ok(!instances.has(canvas), 'canvas initialized twice');
            this.config = config;
            instances.set(canvas, this);
            this.destroy = () => instances.delete(canvas);
        }
    }
    vm.runInNewContext(readFileSync(`${__dirname}/../static/js/charts.js`, 'utf8'), { document, Chart });
    return {
        document, instances, counts, progression,
        emit: (name, detail) => listeners.get(name)?.({ detail }),
        setDark: () => { dark = true; }
    };
}

test('charts initialize after settling, once per canvas, without animation', () => {
    const f = fixture();
    assert.equal(f.instances.size, 0);
    f.emit('htmx:afterSwap');
    assert.equal(f.instances.size, 0);
    f.emit('htmx:afterSettle');
    f.emit('DOMContentLoaded');
    f.emit('htmx:afterSettle');
    assert.equal(f.instances.size, 2);
    for (const chart of f.instances.values()) assert.equal(chart.config.options.animation, false);
});

test('repeated replacements release old charts and history restores new canvases', () => {
    const f = fixture();
    f.emit('DOMContentLoaded');
    for (let i = 0; i < 40; i++) {
        f.emit('htmx:beforeCleanupElement', { elt: { tagName: 'DIV' } });
        for (const canvas of f.document.canvases) f.emit('htmx:beforeCleanupElement', { elt: canvas });
        assert.equal(f.instances.size, 0);
        f.document.canvases = [f.counts(), f.progression()];
        f.emit(i % 2 ? 'htmx:afterSettle' : 'htmx:historyRestore');
        assert.equal(f.instances.size, 2);
    }
});

test('null datasets do not crash and missing progression points retain their meaning', () => {
    const f = fixture();
    f.emit('DOMContentLoaded');
    const line = f.instances.get(f.document.canvases[1]).config;
    assert.equal(JSON.stringify(line.data.datasets[0].data), '[3,null,2]');
    assert.equal(line.data.datasets[0].spanGaps, true);
    for (const canvas of f.document.canvases) f.emit('htmx:beforeCleanupElement', { elt: canvas });
    f.document.canvases[0].dataset.datasets = 'null';
    f.document.canvases[1].dataset.progression = 'null';
    for (const canvas of f.document.canvases) canvas.dataset.labels = 'null';
    f.emit('htmx:afterSettle');
    for (const chart of f.instances.values()) {
        assert.equal(chart.config.data.datasets.length, 0);
        assert.equal(chart.config.data.labels.length, 0);
    }
});

test('theme changes replace existing charts with matching label colors', () => {
    const f = fixture();
    f.emit('DOMContentLoaded');
    const previous = [...f.instances.values()];
    f.setDark();
    f.emit('boulderlog:themeChanged');
    assert.equal(f.instances.size, 2);
    for (const chart of f.instances.values()) {
        assert.ok(!previous.includes(chart));
        assert.equal(chart.config.options.scales.x.ticks.color, '#fff');
    }
});
