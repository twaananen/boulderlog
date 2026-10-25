(() => {
    function gradeCounts(canvas) {
        const labels = JSON.parse(canvas.dataset.labels) ?? [];
        const datasets = JSON.parse(canvas.dataset.datasets) ?? {};
        const isDarkMode = document.documentElement.classList.contains('dark');
        const colors = {
            Topped: { light: 'rgba(54, 162, 235, 0.8)', dark: 'rgba(54, 162, 235, 0.5)' },
            Untopped: { light: 'rgba(255, 99, 132, 0.8)', dark: 'rgba(255, 99, 132, 0.5)' },
            Flashed: { light: 'rgba(255, 165, 0, 0.8)', dark: 'rgba(255, 165, 0, 0.5)' },
            New: { light: 'rgba(76, 175, 80, 0.8)', dark: 'rgba(76, 175, 80, 0.5)' }
        };

        new Chart(canvas, {
            type: 'bar',
            data: {
                labels,
                datasets: Object.keys(datasets).map(key => ({
                    label: key,
                    data: datasets[key],
                    backgroundColor: isDarkMode ? colors[key].dark : colors[key].light,
                    borderColor: isDarkMode ? colors[key].dark.replace('0.5', '1') : colors[key].light.replace('0.8', '1'),
                    borderWidth: 1
                }))
            },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: false,
                scales: {
                    y: {
                        beginAtZero: true,
                        ticks: {
                            color: isDarkMode ? '#fff' : '#333',
                            stepSize: 1,
                            precision: 0
                        },
                        grid: {
                            color: isDarkMode ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)'
                        }
                    },
                    x: {
                        ticks: {
                            color: isDarkMode ? '#fff' : '#333',
                            maxRotation: 90,
                            minRotation: 0
                        },
                        grid: { display: false }
                    }
                },
                plugins: {
                    legend: {
                        labels: { color: isDarkMode ? '#fff' : '#333' },
                        position: 'bottom'
                    }
                }
            }
        });
    }

    function difficultyProgression(canvas) {
        const progressionData = JSON.parse(canvas.dataset.progression) ?? {};
        const grades = Object.keys(progressionData);
        const labels = JSON.parse(canvas.dataset.labels) ?? [];
        const isDarkMode = document.documentElement.classList.contains('dark');
        const textColor = isDarkMode ? '#fff' : '#333';
        const datasets = grades.map((grade, index) => {
            const values = progressionData[grade].map(point => point.Value);
            const lastTwoValues = values.filter(v => v > 0).slice(-2);
            const shouldHide = lastTwoValues.length === 2 && lastTwoValues.every(val => val === 1);

            return {
                label: grade,
                data: values.map(value => value > 0 ? value : null),
                borderColor: `hsl(${(index * 360 / grades.length) % 360}, 70%, 50%)`,
                backgroundColor: `hsla(${(index * 360 / grades.length) % 360}, 70%, 50%, 0.2)`,
                tension: 0.3,
                pointRadius: 4,
                pointHoverRadius: 6,
                borderWidth: 2,
                hidden: shouldHide,
                spanGaps: true
            };
        });

        new Chart(canvas, {
            type: 'line',
            data: { labels, datasets },
            options: {
                responsive: true,
                maintainAspectRatio: false,
                animation: false,
                interaction: { mode: 'nearest', axis: 'x', intersect: false },
                scales: {
                    x: {
                        ticks: { color: textColor, maxRotation: 45, minRotation: 45 },
                        grid: { color: isDarkMode ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)' }
                    },
                    y: {
                        min: 1,
                        max: 8,
                        ticks: {
                            color: textColor,
                            stepSize: 1,
                            autoSkip: false,
                            callback: function(value) {
                                return Number.isInteger(value) ? value : value.toFixed(1);
                            }
                        },
                        grid: { color: isDarkMode ? 'rgba(255, 255, 255, 0.1)' : 'rgba(0, 0, 0, 0.1)' }
                    }
                },
                plugins: {
                    legend: {
                        position: 'bottom',
                        labels: { color: textColor, padding: 20, usePointStyle: true, pointStyle: 'circle' }
                    },
                    tooltip: {
                        callbacks: {
                            label: function(context) {
                                return `${context.dataset.label}: ${context.parsed.y.toFixed(1)}`;
                            }
                        }
                    }
                }
            }
        });
    }

    function initializeCharts() {
        document.querySelectorAll('canvas[data-chart]').forEach(canvas => {
            if (Chart.getChart(canvas)) return;
            if (canvas.dataset.chart === 'counts') gradeCounts(canvas);
            if (canvas.dataset.chart === 'progression') difficultyProgression(canvas);
        });
    }

    document.addEventListener('DOMContentLoaded', initializeCharts);
    // HTMX settles canvas dimensions after swapping; initializing earlier clears the drawing.
    document.addEventListener('htmx:afterSettle', initializeCharts);
    document.addEventListener('htmx:historyRestore', initializeCharts);
    document.addEventListener('htmx:beforeCleanupElement', event => {
        if (event.detail.elt.tagName === 'CANVAS') Chart.getChart(event.detail.elt)?.destroy();
    });
    document.addEventListener('boulderlog:themeChanged', () => {
        document.querySelectorAll('canvas[data-chart]').forEach(canvas => Chart.getChart(canvas)?.destroy());
        initializeCharts();
    });
})();
