[![Contributors][contributors-shield]][contributors-url]
[![Language][language-shield]][language-url]
[![Issues][issues-shield]][issues-url]
[![GPL-3.0 License][license-shield]][license-url]
[![Watchers][watchers-shield]][watchers-url]

<!-- PROJECT LOGO -->
<br />
<p align="center">

  <h3 align="center">Granular Data Gatherer (gdg)</h3>

  <p align="center">
    Collects Granular OS Metrics for Troubleshooting
    <br />
    <a href="https://github.com/rfparedes/gdg/issues">Report Bug</a>
    ·
    <a href="https://github.com/rfparedes/gdg/issues">Request Feature</a>
  </p>
</p>

<!-- TABLE OF CONTENTS -->
<details open="open">
  <summary><h2 style="display: inline-block">Table of Contents</h2></summary>
  <ol>
    <li>
      <a href="#about-the-project">About The Project</a>
    </li>
    <li>
      <a href="#getting-started">Getting Started</a>
      <ul>
        <li><a href="#prerequisites">Prerequisites</a></li>
        <li><a href="#installation">Installation</a></li>
      </ul>
    </li>
    <li><a href="#technical-details">Technical Details</a></li>
    <li><a href="#usage">Usage</a></li>
    <li><a href="#build-it-yourself">Build It Yourself</a></li>
    <li><a href="#roadmap">Roadmap</a></li>
    <li><a href="#contributing">Contributing</a></li>
    <li><a href="#license">License</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About The Project

gdg or Granular Data Gatherer was developed in Go to fill the missing gap in the availability of an easy and open all-in-one tool to collect OS metrics for troubleshooting.  OSWatcher and nmon cannot be the only viable options.

<!-- GETTING STARTED -->
## Getting Started

To get a local copy up and running follow these simple steps.

### Prerequisites

* a server, instance, VM running a systemd-enabled Linux distribution

### Installation

Download the binary from Releases
<https://github.com/rfparedes/gdg/releases/download/v0.9.0/gdg>

Start it

```sh
sudo ./gdg -start
```

## Technical Details

* gdg uses standard Linux utilities to perform its work, including:

  * iostat
  * top
  * mpstat
  * vmstat
  * ss
  * nstat
  * ps
  * nfsiostat
  * ethtool
  * ip
  * pidstat

* gdg will detect which utilities are available and only use those.  In advance, you can install any of the utilities above anytime before or after setup.

* gdg will create a configuration file and data directory in the same directory where the gdg binary resides. e.g. If you download the binary to `/usr/local/` this directory is where metric data and config file will be stored.

* gdg uses a systemd timer so there is no running daemon

* gdg only installs a systemd service and systemd timer on `-start` outside of the working directory where the gdg binary resides

* gdg removes the systemd service and systemd timer on `-stop`.  The working directory is untouched.

## Usage

### To start collection in 30s intervals, run

```sh
sudo ./gdg -t 30 -start
```

### To stop collection, run

```sh
sudo ./gdg -stop
```

### To see the current status of gdg including start/stop status, version, interval, data location, and current size of collected data, run

```sh
sudo ./gdg -status
```

e.g.

```
VERSION: gdg-0.9.0
STATUS: started
INTERVAL: 30s
DATA LOCATION: /usr/local/gdg/gdg-data/
CURRENT DATA SIZE: 79MB
```

### For help

```sh
./gdg -h
```

## Build it yourself

* You'll need a go compiler installed

Clone it

```sh
git clone https://github.com/rfparedes/gdg.git
```

Build it

```sh
cd gdg
go build -o gdg
```

Start it

```sh
sudo ./gdg -start
```
<!-- ROADMAP -->
## Roadmap

See the [open issues](https://github.com/rfparedes/gdg/issues) for a list of proposed features (and known issues).

<!-- CONTRIBUTING -->
## Contributing

Contributions are what make the open source community such an amazing place to be learn, inspire, and create. Any contributions you make are **greatly appreciated**.

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

<!-- LICENSE -->
## License

Distributed under the GPL-3.0 License. See `LICENSE` for more information.

<!-- MARKDOWN LINKS & IMAGES -->
<!-- https://www.markdownguide.org/basic-syntax/#reference-style-links -->
[contributors-shield]: https://img.shields.io/github/contributors/rfparedes/gdg?color=%20%2330BA78
[contributors-url]: https://github.com/rfparedes/repo/graphs/contributors
[language-shield]: https://img.shields.io/github/languages/top/rfparedes/gdg?color=%20%2330BA78
[language-url]: https://github.com/rfparedes/gdg/search?l=go
[watchers-shield]: https://img.shields.io/github/watchers/rfparedes/gdg?color=%20%2330BA78&style=social
[watchers-url]:https://github.com/rfparedes/gdg/watchers
[issues-shield]: https://img.shields.io/github/issues/rfparedes/gdg?color=%20%2330BA78
[issues-url]: https://github.com/rfparedes/gdg/issues
[license-shield]: https://img.shields.io/github/license/rfparedes/gdg?color=%20%2330BA78
[license-url]: https://github.com/rfparedes/gdg/blob/main/LICENSE
