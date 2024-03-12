import './App.css';

import DatePicker from 'react-datepicker'
import "react-datepicker/dist/react-datepicker.css";
import { useState } from 'react';

import { Line } from "react-chartjs-2";
import { Chart } from "chart.js/auto";
import { CategoryScale } from "chart.js"

Chart.register(CategoryScale);

function fetchData(endpoint, key) {
  let res = localStorage.getItem(key);
  if (res !== null) {
    let parsed = JSON.parse(res);

    if (Date.now() > parsed.expires) {
      localStorage.removeItem(key);
      return fetchData(endpoint, key);
    }
    return Promise.resolve(parsed.response.records);
  }

  return fetch(endpoint + key)
    .then((response) => response.json())
    .then((data) => {
      localStorage.setItem(key, JSON.stringify({
        expires: Date.now() + 500000,
        response: data
      }));
      return data.records;
    })
    .catch((err) => {
      console.log(err);
      return [];
    });
}

function parseResponse(res) {
  let parseUnixtime = (ts) => {
    let options = {
      year: "numeric",
      month: "long",
      day: "numeric",
      hour: "numeric"
    };
    let date = new Date(ts);

    return date.toLocaleDateString("en-US", options);
  };

  return res.map((d) => {
    return {
      x: parseUnixtime(d.unixtime),
      y: d.posts
    }
  });
}

/*
 * [{ 
 *    unixtime: uint,
 *    posts: uint
 * }]
 *
 * */
function extractRange(array, lb, ub) {
  if (array.length === 0) {
    return array;
  }
  let k = 0;
  let len = array.length;
  for (var z = len-1; z > 0; z = Math.floor(z/2)) {
    while (k+z < len && array[k+z].unixtime <= lb) {
      k += z;
    }
  }
  console.log(lb, ub, k, array);
  let result = [];
  for (; k < len && array[k].unixtime <= ub; k++) {
    if (array[k].unixtime < lb) {
      continue;
    }
    result.push(array[k]);
  }
  return result;
}

function buildData(array) {
  return {
    labels: array.map((data) => data.x),
    datasets: [
      {
        label: "Посты",
        data: array.map((data) => data.y),
        backgroundColor: [
          "rgba(75,192,192,1)",
          "#ecf0f1",
          "#50AF95",
          "#f3ba2f",
          "#2a71d0"
        ],
        borderColor: "#494fc6",
        borderWidth: 2
      }
    ]
  };
}

function LineChart(props) {
  return ( 
    <Line
      data={props.chartData}
      options={{
        plugins: {
          title: {
            display: false,
            text: props.title
          },
          legend: {
            display: false
          },
        },
        layout: {
          padding: 20
        },
        scales: {
          y: {
            min: 0
          },
        }
      }}
    />
  );
}

function App() {
  const [chartData, nextChartData] = useState(buildData([
    //{x: 1, y: 1.1},
  ]))

  const [dateRange, setDateRange] = useState([null, null]);
  const [startDate, endDate] = dateRange;

  const [board, nextBoard] = useState("2ch.hk/b");
  const handleBoardState = (e) => {
    nextBoard(e.target.value);
    setDateRange([null, null]);
  };

  const callback = (key, lb, ub) => {
    if (lb === null || ub === null) {
      console.log("left or right bound is null");
      return;
    }
    fetchData("http://localhost:8080/api/stats?board=", key)
      .then((data) => {
        console.log(data);
        return extractRange(data, lb.getTime(), ub.getTime());
      })
      .then((ranged) => {
        let parsed = parseResponse(ranged);
        nextChartData(buildData(parsed));
      })
      .catch((e) => {
        console.log(e);
      });
  };

  return (
    <div>
      <div class="navbar">
        [<a href="/">home</a>]
        [info]
      </div>
      <div class="header">
        kuklobund intelligence agency
      </div>
      <hr />
      <div class="horizontal-wrapper">
        <div class="settings">
          <div class="settings-header">
            settings
          </div>
          <DatePicker
            selectsRange={true}
            startDate={startDate}
            endDate={endDate}
            onChange={(update) => {
                setDateRange(update);
                callback(board, update[0], update[1]);
            }}
            placeholderText='выбрать период'
            className='picker'
          />
          <select value={board} onChange={handleBoardState}>
            <option value="2ch.hk/b">2ch.hk/b</option>
            <option value="2ch.hk/po">2ch.hk/po</option>
            <option value="2ch.hk/vg">2ch.hk/vg</option>  
          </select>
        </div>
        <div class="chart-container">
          <div class="chart-header">
            posting stats
          </div>
          <LineChart chartData={chartData} />
        </div>
      </div>
    </div>
  );
}

export default App;
