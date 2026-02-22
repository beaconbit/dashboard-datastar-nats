#!/usr/bin/env python3
import re

with open('components/dashboard.templ', 'r') as f:
    lines = f.readlines()

layout = None  # 'tv' or 'mobile'
quarter = None  # 'A','B','C','D'
counters = {
    'tv': {'A': 0, 'B': 0, 'C_wait': 0, 'C_wasted': 0, 'D': 0},
    'mobile': {'A': 0, 'B': 0, 'C_wait': 0, 'C_wasted': 0, 'D': 0}
}
wasted_keys = ['hour4', 'hour3', 'hour2', 'hour1', 'current']
wasted_idx = 0

output = []
i = 0
while i < len(lines):
    line = lines[i]
    stripped = line.strip()
    
    # Detect layout
    if '<div class="tv"' in line:
        layout = 'tv'
    elif '<div class="mobile"' in line:
        layout = 'mobile'
    
    # Detect quarter based on surrounding content
    if 'quarter-a-content' in line:
        quarter = 'A'
        counters[layout]['A'] = 0  # reset column count for this quarter
    elif 'quarter-b-content' in line:
        quarter = 'B'
        counters[layout]['B'] = 0
    elif 'quarter-c-content' in line:
        quarter = 'C'
        counters[layout]['C_wait'] = 0
        counters[layout]['C_wasted'] = 0
        wasted_idx = 0
    elif 'quota-columns' in line:
        quarter = 'D'
        counters[layout]['D'] = 0
    
    # Process dynamic elements
    if layout and quarter:
        # Square grid for quarter A or B
        if 'class="square-grid"' in line:
            if quarter in ('A', 'B'):
                col_type = quarter
                counters[layout][col_type] += 1
                col_idx = counters[layout][col_type]
                new_id = f'{layout}-quarter{col_type}-col{col_idx}'
                if 'id=' in line:
                    line = re.sub(r'id="[^"]*"', f'id="{new_id}"', line)
                else:
                    line = line.replace('class="square-grid"', f'id="{new_id}" class="square-grid"')
        # Wait time square for quarter C
        elif 'class="wait-time-square"' in line:
            if quarter == 'C':
                counters[layout]['C_wait'] += 1
                wait_idx = counters[layout]['C_wait']
                new_id = f'{layout}-quarterC-cbw{wait_idx}-waitTime'
                if 'id=' in line:
                    line = re.sub(r'id="[^"]*"', f'id="{new_id}"', line)
                else:
                    line = line.replace('class="wait-time-square"', f'id="{new_id}" class="wait-time-square"')
        # Wasted minute value for quarter C
        elif 'class="wasted-minute-value"' in line:
            if quarter == 'C' and wasted_idx < len(wasted_keys):
                key = wasted_keys[wasted_idx]
                new_id = f'{layout}-quarterC-wastedMinutes-{key}'
                if 'id=' in line:
                    line = re.sub(r'id="[^"]*"', f'id="{new_id}"', line)
                else:
                    line = line.replace('class="wasted-minute-value"', f'id="{new_id}" class="wasted-minute-value"')
                wasted_idx += 1
        # Quota column for quarter D
        elif 'class="quota-column"' in line:
            if quarter == 'D':
                counters[layout]['D'] += 1
                col_idx = counters[layout]['D']
                new_id = f'{layout}-quarterD-col{col_idx}'
                if 'id=' in line:
                    line = re.sub(r'id="[^"]*"', f'id="{new_id}"', line)
                else:
                    line = line.replace('class="quota-column"', f'id="{new_id}" class="quota-column"')
    
    output.append(line)
    i += 1

new_content = ''.join(output)

with open('components/dashboard.templ', 'w') as f:
    f.write(new_content)

print('Added missing IDs and fixed existing IDs.')