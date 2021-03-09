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
    <li><a href="#validated-distributions">Validated Distributions</a></li>
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

Download the binary from Releases (<https://github.com/rfparedes/gdg/releases/latest/download/gdg>) to `/usr/local/sbin` on the server and run:

```sh
sudo chmod +x /usr/local/sbin/gdg
```

Start it

```sh
sudo /usr/local/sbin/gdg -start
```

Check Status Anytime

```sh
/usr/local/sbin/gdg -status
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
  * rtmon

* gdg will detect which utilities are available and only use those installed.  In advance, you can install any of the utilities above anytime before or after setup. Most of these utilities are located in only five different packages. On most distributions, sysstat package contains (`iostat`, `mpstat`, `pidstat`), nfs-common or nfs-client package contains (`nfsiostat`), procps package contains (`top`, `vmstat`, `ps`), iproute2 package contains (`ss`, `nstat`, `ip`, `rtmon`) and ethtool contains (`ethtool`).

* gdg will create a configuration file in `/etc/gdg.cfg` and a data directory in `/var/log/gdg-data`.

* gdg uses a systemd timer so there is no running daemon.

* gdg installs a systemd service and systemd timer on `-start`.

* gdg removes the systemd service and systemd timer on `-stop`.  All other files are untouched.

* gdg collects data in the `/var/log/gdg-data` directory.  The children below this directory are named after the utility (e.g. `iostat`) which collected the data.  Below this directory are .dat (e.g. `meminfo_21.03.07.2300.dat`) files named after the following format (`utility_YY.MM.DD.HH00.dat`). The .dat files contain at maximum, one hour worth of data.

* To easily search down chronologically through the data collected in the .dat file, use the search string `zzz`.

## Usage

### To start collection in 30s intervals, run

```sh
sudo /usr/local/sbin/gdg -t 30 -start
```

### To stop collection, run

```sh
sudo /usr/local/sbin/gdg -stop
```

### To see the data collected

```sh
cd /var/log/gdg-data
```

### To see the current status of gdg including start/stop status, version, interval, data location, and current size of collected data, run

```sh
/usr/local/sbin/gdg -status
```

e.g.

```
VERSION: gdg-0.9.0
STATUS: started
INTERVAL: 30s
DATA LOCATION: /var/log/gdg-data/
CONFIG LOCATION: /etc/gdg.cfg
CURRENT DATA SIZE: 79MB
```

### If you want to change the interval (-t) or after installing additional supported utilities, run

```sh
sudo /usr/local/sbin/gdg -reload
```

### For help

```sh
/usr/local/sbin/gdg -h
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

Move it

```sh
mv gdg /usr/local/sbin
sudo chmod +x /usr/local/sbin/gdg
```

Start it

```sh
sudo /usr/local/sbin/gdg -start
```

## Validated Distributions

gdg has been validated on:

* SLE-12 (SLES or SLES-SAP 12 all SPs)
* SLE-15 (SLES or SLES-SAP 15 all SPs)
* openSUSE Leap 12/15
* Debian 9
* Debian 10
* RHEL7
* RHEL8

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
[contributors-url]: https://github.com/rfparedes/gdg/graphs/contributors
[language-shield]: https://img.shields.io/github/languages/top/rfparedes/gdg?color=%20%2330BA78
[language-url]: https://github.com/rfparedes/gdg/search?l=go
[watchers-shield]: https://img.shields.io/github/watchers/rfparedes/gdg?color=%20%2330BA78&style=social
[watchers-url]:https://github.com/rfparedes/gdg/watchers
[issues-shield]: https://img.shields.io/github/issues/rfparedes/gdg?color=%20%2330BA78
[issues-url]: https://github.com/rfparedes/gdg/issues
[license-shield]: https://img.shields.io/github/license/rfparedes/gdg?color=%20%2330BA78
[license-url]: https://github.com/rfparedes/gdg/blob/main/LICENSE
