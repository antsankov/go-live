import matplotlib
matplotlib.use('Agg')
import matplotlib.pyplot as plt
import numpy as np

servers = ['go-live', 'caddy', 'miniserve', 'live-server']
colors = ['#2196F3', '#FF9800', '#4CAF50', '#F44336']

data = {
    '1KB': {
        1:  [4888, 3333, 3616, 2147],
        10: [18484, 8169, 14944, 5499],
        50: [22421, 7734, 13185, 4716],
    },
    '50KB': {
        1:  [3598, 2423, 3134, 2556],
        10: [12681, 7854, 12180, 9931],
        50: [12881, 7953, 12003, 10548],
    },
    '1MB': {
        1:  [1006, 1016, 702, 656],
        10: [1642, 1549, 1288, 1165],
        50: [1615, 1590, 1292, 1306],
    },
}

fig, axes = plt.subplots(1, 3, figsize=(18, 6))
fig.suptitle('File Server Benchmark — Requests/sec (higher is better)', fontsize=16, fontweight='bold', y=1.02)

concurrencies = [1, 10, 50]
x = np.arange(len(concurrencies))
width = 0.18

for idx, (file_size, conc_data) in enumerate(data.items()):
    ax = axes[idx]

    for i, (server, color) in enumerate(zip(servers, colors)):
        values = [conc_data[c][i] for c in concurrencies]
        offset = (i - 1.5) * width
        bars = ax.bar(x + offset, values, width, label=server, color=color, edgecolor='white', linewidth=0.5)

        for bar, val in zip(bars, values):
            ax.text(bar.get_x() + bar.get_width()/2., bar.get_height() + max(values)*0.01,
                    f'{val:,}', ha='center', va='bottom', fontsize=7, fontweight='bold')

    ax.set_xlabel('Concurrent Clients', fontsize=12)
    ax.set_ylabel('Requests/sec', fontsize=12)
    ax.set_title(f'{file_size} files', fontsize=14, fontweight='bold')
    ax.set_xticks(x)
    ax.set_xticklabels(concurrencies)
    ax.legend(fontsize=9)
    ax.grid(axis='y', alpha=0.3)
    ax.set_axisbelow(True)

plt.tight_layout()
import os
out = os.path.join(os.path.dirname(os.path.abspath(__file__)), 'benchmark.png')
plt.savefig(out, dpi=150, bbox_inches='tight')
print(f'Saved to {out}')
