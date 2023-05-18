<?php

// Make a GET request to the exchange API
$url = 'http://172.18.0.1:6000/exchange';
$response = file_get_contents($url);

// Decode the JSON response into a PHP array
$data = json_decode($response, true);

// Find the index of the row corresponding to today's date
$today = date('Y-m-d');
$today_index = -1;
for ($i = 0; $i < count($data); $i++) {
  if ($data[$i]['date'] === $today) {
    $today_index = $i;
    break;
  }
}
// Print the row for today's date
if ($today_index !== -1) {
    echo '<h2>Exchange rates against USD for today :</h2>';
    echo '<p>EUR: ' . $data[$today_index]['EUR'] . '</p>';
    echo '</p>GDP: ' . $data[$today_index]['GBP'] . '</p>';
  }
// Print the table headers
echo '<h2>Historical Data :</h2>';
echo '<table>';
echo '<tr><th>Date</th><th>EUR</th><th>GBP</th></tr>';

// Print the history rows
for ($i = 0; $i < count($data); $i++) {
  if ($i !== $today_index) {
    echo '<tr>';
    echo '<td style="border: 1px solid #ddd; padding: 8px;">' . $data[$i]['date'] . '</td>';
    echo '<td style="border: 1px solid #ddd; padding: 8px;">' . $data[$i]['EUR'] . '</td>';
    echo '<td style="border: 1px solid #ddd; padding: 8px;">' . $data[$i]['GBP'] . '</td>';
    echo '</tr>';
  }
}
echo "</table>"



?>
