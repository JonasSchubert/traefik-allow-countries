# Traefik Allow Countries

A [Traefik](https://github.com/containous/traefik) plugin to allow only certain countries and block everything else. Uses the hourly updated country IP blocks from [here](https://github.com/herrbischoff/country-ip-blocks).

## Configuration

Sample configuration in Traefik.

### Configuration as local plugin

traefik.yml

```yaml
log:
  level: INFO
experimental:
  localPlugins:
    allow-countries:
      moduleName: github.com/jonasschubert/traefik-allow-countries
```

dynamic-configuration.yml

```yaml
http:
  middlewares:
    allow-countries-de:
      plugin:
        allow-countries:
          addCountryHeader: true
          allowLocalRequests: true
          cidrFileFolder: /usr/traefik/plugins/cidr
          cidrFileUpdate: true
          countries:
            - DE
          logAllowedRequests: false
          logLocalRequests: false
          silentStartUp: true
```

docker-compose.yml

```yaml
version: "3.7"
services:
  traefik:
    image: traefik
    volumes:
      - /var/run/docker.sock:/var/run/docker.sock
      - /docker/config/traefik/traefik.yml:/etc/traefik/traefik.yml
      - /docker/config/traefik/dynamic-configuration.yml:/etc/traefik/dynamic-configuration.yml
      - /plugin/allow-countries:/plugins-local/src/github.com/jonasschubert/traefik-allow-countries/
    ports:
      - "80:80"
  hello:
    image: containous/whoami
    labels:
      - traefik.enable=true
      - traefik.http.routers.hello.entrypoints=http
      - traefik.http.routers.hello.rule=Host(`hello.localhost`)
      - traefik.http.services.hello.loadbalancer.server.port=80
      - traefik.http.routers.hello.middlewares=my-plugin@file
```

## Sample configuration

- `addCountryHeader`:  If set to `true`, adds the X-IPCountry header to the HTTP request header. The header contains the two letter country code returned by cache or API request.
- `allowLocalRequests`: If set to true, will not block request from [Private IP Ranges](https://de.wikipedia.org/wiki/Private_IP-Adresse)
- `countries`: list of allowed countries
- `logLocalRequests`: If set to true, will log every connection from any IP in the private IP range

```yaml
my-allow-countries:
  plugin:
    allow-countries:
      addCountryHeader: true
      allowLocalRequests: true
      cidrFileFolder: /usr/traefik/plugins/cidr
      cidrFileUpdate: true
      countries:
        - AF # Afghanistan
        - AL # Albania
        - DZ # Algeria
        - AS # American Samoa
        - AD # Andorra
        - AO # Angola
        - AI # Anguilla
        - AQ # Antarctica
        - AG # Antigua and Barbuda
        - AR # Argentina
        - AM # Armenia
        - AW # Aruba
        - AU # Australia
        - AT # Austria
        - AZ # Azerbaijan
        - BS # Bahamas (the)
        - BH # Bahrain
        - BD # Bangladesh
        - BB # Barbados
        - BY # Belarus
        - BE # Belgium
        - BZ # Belize
        - BJ # Benin
        - BM # Bermuda
        - BT # Bhutan
        - BO # Bolivia (Plurinational State of)
        - BQ # Bonaire, Sint Eustatius and Saba
        - BA # Bosnia and Herzegovina
        - BW # Botswana
        - BV # Bouvet Island
        - BR # Brazil
        - IO # British Indian Ocean Territory (the)
        - BN # Brunei Darussalam
        - BG # Bulgaria
        - BF # Burkina Faso
        - BI # Burundi
        - CV # Cabo Verde
        - KH # Cambodia
        - CM # Cameroon
        - CA # Canada
        - KY # Cayman Islands (the)
        - CF # Central African Republic (the)
        - TD # Chad
        - CL # Chile
        - CN # China
        - CX # Christmas Island
        - CC # Cocos (Keeling) Islands (the)
        - CO # Colombia
        - KM # Comoros (the)
        - CD # Congo (the Democratic Republic of the)
        - CG # Congo (the)
        - CK # Cook Islands (the)
        - CR # Costa Rica
        - HR # Croatia
        - CU # Cuba
        - CW # Curaçao
        - CY # Cyprus
        - CZ # Czechia
        - CI # Côte d'Ivoire
        - DK # Denmark
        - DJ # Djibouti
        - DM # Dominica
        - DO # Dominican Republic (the)
        - EC # Ecuador
        - EG # Egypt
        - SV # El Salvador
        - GQ # Equatorial Guinea
        - ER # Eritrea
        - EE # Estonia
        - SZ # Eswatini
        - ET # Ethiopia
        - FK # Falkland Islands (the) [Malvinas]
        - FO # Faroe Islands (the)
        - FJ # Fiji
        - FI # Finland
        - FR # France
        - GF # French Guiana
        - PF # French Polynesia
        - TF # French Southern Territories (the)
        - GA # Gabon
        - GM # Gambia (the)
        - GE # Georgia
        - DE # Germany
        - GH # Ghana
        - GI # Gibraltar
        - GR # Greece
        - GL # Greenland
        - GD # Grenada
        - GP # Guadeloupe
        - GU # Guam
        - GT # Guatemala
        - GG # Guernsey
        - GN # Guinea
        - GW # Guinea-Bissau
        - GY # Guyana
        - HT # Haiti
        - HM # Heard Island and McDonald Islands
        - VA # Holy See (the)
        - HN # Honduras
        - HK # Hong Kong
        - HU # Hungary
        - IS # Iceland
        - IN # India
        - ID # Indonesia
        - IR # Iran (Islamic Republic of)
        - IQ # Iraq
        - IE # Ireland
        - IM # Isle of Man
        - IL # Israel
        - IT # Italy
        - JM # Jamaica
        - JP # Japan
        - JE # Jersey
        - JO # Jordan
        - KZ # Kazakhstan
        - KE # Kenya
        - KI # Kiribati
        - KP # Korea (the Democratic People's Republic of)
        - KR # Korea (the Republic of)
        - KW # Kuwait
        - KG # Kyrgyzstan
        - LA # Lao People's Democratic Republic (the)
        - LV # Latvia
        - LB # Lebanon
        - LS # Lesotho
        - LR # Liberia
        - LY # Libya
        - LI # Liechtenstein
        - LT # Lithuania
        - LU # Luxembourg
        - MO # Macao
        - MG # Madagascar
        - MW # Malawi
        - MY # Malaysia
        - MV # Maldives
        - ML # Mali
        - MT # Malta
        - MH # Marshall Islands (the)
        - MQ # Martinique
        - MR # Mauritania
        - MU # Mauritius
        - YT # Mayotte
        - MX # Mexico
        - FM # Micronesia (Federated States of)
        - MD # Moldova (the Republic of)
        - MC # Monaco
        - MN # Mongolia
        - ME # Montenegro
        - MS # Montserrat
        - MA # Morocco
        - MZ # Mozambique
        - MM # Myanmar
        - NA # Namibia
        - NR # Nauru
        - NP # Nepal
        - NL # Netherlands (the)
        - NC # New Caledonia
        - NZ # New Zealand
        - NI # Nicaragua
        - NE # Niger (the)
        - NG # Nigeria
        - NU # Niue
        - NF # Norfolk Island
        - MP # Northern Mariana Islands (the)
        - NO # Norway
        - OM # Oman
        - PK # Pakistan
        - PW # Palau
        - PS # Palestine, State of
        - PA # Panama
        - PG # Papua New Guinea
        - PY # Paraguay
        - PE # Peru
        - PH # Philippines (the)
        - PN # Pitcairn
        - PL # Poland
        - PT # Portugal
        - PR # Puerto Rico
        - QA # Qatar
        - MK # Republic of North Macedonia
        - RO # Romania
        - RU # Russian Federation (the)
        - RW # Rwanda
        - RE # Réunion
        - BL # Saint Barthélemy
        - SH # Saint Helena, Ascension and Tristan da Cunha
        - KN # Saint Kitts and Nevis
        - LC # Saint Lucia
        - MF # Saint Martin (French part)
        - PM # Saint Pierre and Miquelon
        - VC # Saint Vincent and the Grenadines
        - WS # Samoa
        - SM # San Marino
        - ST # Sao Tome and Principe
        - SA # Saudi Arabia
        - SN # Senegal
        - RS # Serbia
        - SC # Seychelles
        - SL # Sierra Leone
        - SG # Singapore
        - SX # Sint Maarten (Dutch part)
        - SK # Slovakia
        - SI # Slovenia
        - SB # Solomon Islands
        - SO # Somalia
        - ZA # South Africa
        - GS # South Georgia and the South Sandwich Islands
        - SS # South Sudan
        - ES # Spain
        - LK # Sri Lanka
        - SD # Sudan (the)
        - SR # Suriname
        - SJ # Svalbard and Jan Mayen
        - SE # Sweden
        - CH # Switzerland
        - SY # Syrian Arab Republic
        - TW # Taiwan (Province of China)
        - TJ # Tajikistan
        - TZ # Tanzania, United Republic of
        - TH # Thailand
        - TL # Timor-Leste
        - TG # Togo
        - TK # Tokelau
        - TO # Tonga
        - TT # Trinidad and Tobago
        - TN # Tunisia
        - TR # Turkey
        - TM # Turkmenistan
        - TC # Turks and Caicos Islands (the)
        - TV # Tuvalu
        - UG # Uganda
        - UA # Ukraine
        - AE # United Arab Emirates (the)
        - GB # United Kingdom of Great Britain and Northern Ireland (the)
        - UM # United States Minor Outlying Islands (the)
        - US # United States of America (the)
        - UY # Uruguay
        - UZ # Uzbekistan
        - VU # Vanuatu
        - VE # Venezuela (Bolivarian Republic of)
        - VN # Viet Nam
        - VG # Virgin Islands (British)
        - VI # Virgin Islands (U.S.)
        - WF # Wallis and Futuna
        - EH # Western Sahara
        - YE # Yemen
        - ZM # Zambia
        - ZW # Zimbabwe
        - AX # Åland Islands
      logAllowedRequests: false
      logLocalRequests: true
      silentStartUp: false
```

## Configuration options

### Allow local requests: `allowLocalRequests`

If set to true, will not block request from [Private IP Ranges](https://en.wikipedia.org/wiki/Private_network).

Defaults to `false`.

### CIDR file folder `cidrFileFolder`

The (mounted) folder with the sub folders `ipv4` and `ipv6` containing the `CIDR` files for each country.

Recommandation is to clone [this project](https://github.com/herrbischoff/country-ip-blocks) and mount it as well as update it in a cron job.

### CIDR file update `cidrFileUpdate`

Enables the hourly update for the CIDR files. Recommended to set to `true`.

Defaults to `true`.

### Countries `countries`

A list of country codes from which connections to the service should be allowed.

### Log allowed requests `logAllowedRequests`

If set to true, will show a log message with the IP and the country of origin if a request is allowed.

Defaults to `false`.

### Log details `logDetails`

If set to true, will show a log message with method calls and more.

Defaults to `true`.

### Log local requests: `logLocalRequests`

If set to true, will show a log message when some one accesses the service over a private ip address.

Defaults to `true`.

## Contributors

| [<img alt="JonasSchubert" src="https://secure.gravatar.com/avatar/835215bfb654d58acb595c64f107d052?s=180&d=identicon" width="117"/>](https://code.schubert.zone/jonas-schubert) |
| :---------------------------------------------------------------------------------------------------------------------------------------: |
| [Jonas Schubert](https://code.schubert.zone/jonas-schubert) |

## License

traefik-allow-countries is distributed under the MIT license. [See LICENSE](LICENSE) for details.

```
MIT License

Copyright (c) 2022-2023 Jonas Schubert

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
```
