# goldfinger

Retrieves US Savings Bond information from the US Treasury and provides some basic information about them:

- Total Value
- Total Interest Earned
- Total Original Purcahse Price
- Which bonds still have to mature and when they will

## Configuration

Copy the `config.example.yaml` to `config.yaml` and populate with your information from your bonds

```yaml
bonds:
  - denomination: 50
    serial: L123456789EE
    issue_date: "11/1990"
    series: "EE"
  - denomination: 100
    serial: C123456789EE
    issue_date: "07/1994"
    series: "EE"
```

## Usage

```console
goldfinger -h
```

```console
Usage of ./goldfinger:
  -config string
    	path to config file (default "./config.yaml")
```

## Example

```console
goldfinger -config config.example.yaml
2024/04/02 21:35:36 Reading configuration from configFile=config.example.yaml
2024/04/02 21:35:36 Loaded 2 bonds
2024/04/02 21:35:36 Retrieving bond values
2024/04/02 21:35:37 Total Value = $266.16
2024/04/02 21:35:37 Total Interest = $191.16
2024/04/02 21:35:37 Total Original Purchase Price = $75.00
2024/04/02 21:35:37 Found 1 bonds still to mature
2024/04/02 21:35:37 Bond C123456789EE will mature on 07/2024
```
