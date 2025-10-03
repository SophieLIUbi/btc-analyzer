# Bitcoin Market Analyzer
A comprehensive Bitcoin market analysis tool written in Go that fetches real-time data, performs technical analysis, and generates detailed reports.

## 🚀 Features
-**Multiple Data Sources**: API, CSV, JSON, or sample data  
-**Technical Indicators:** RSI, MACD, Bollinger Bands, Moving Averages  
-**Risk Analysis:** Volatility, Sharpe Ratio, Maximum Drawdown  
-**Pattern Detection:** Support/resistance, trend analysis, volume patterns  
-**Report Generation**: HTML, JSON, and CSV exports  
-**Trading Signals:** Automated buy/sell/hold recommendations  
## 🛠 Installation
**Prerequisites**
-Go 1.19 or later  
-Internet connection (for API data)  
**Setup**

`mkdir btc-analyzer`  
`cd btc-analyzer`  
`go mod init btc-analyzer`  
`go mod tidy`  
`go build -o btc-analyzer .`  

## 🚀 Quick Start

`# Generate sample data and analyze`  
`go run . -source=sample -days=30`  

`# Fetch real Bitcoin data from API`  
`go run . -source=api -days=7`  

`# Open HTML report (Windows)`  
`start output\btc_analysis_report.html`  

`# Open HTML report (Mac/Linux)`  
`open output/btc_analysis_report.html`  

## 📖 Usage Examples
**Basic Usage**  
`# Real-time API analysis`  
`go run . -source=api -days=30`  

`# Sample data for testing`  
`go run . -source=sample -days=60`  

`# Analyze CSV file`  
`go run . -source=csv -csv=./data/prices.csv`  

`# Custom output directory`  
`go run . -source=api -days=14 -output=./reports`  


**Report Generation**

`# Generate all reports`  
`go run . -source=sample -days=30 -verbose`  

`# Only HTML report`  
`go run . -source=api -days=7 -json-report=false`  

`# View JSON report`  
`notepad output\btc_analysis_report.json  # Windows`  
`cat output/btc_analysis_report.json      # Linux/Mac`  

## 📁 Project Structure

**btc-analyzer/  
├── main.go                         # Main application  
├── go.mod                          # Dependencies   
├── README.md                       # Documentation  
├── output/                         # Generated reports  
│   ├── btc_analysis_report.html   # HTML report  
│   ├── btc_analysis_report.json   # JSON report  
│   └── btc_data.csv               # Exported data  
└── internal/                      # Source code  
    ├── types/types.go             # Data structures  
    ├── timeseries/timeseries.go   # Time series utils  
    ├── statistics/statistics.go   # Statistical calculations  
    ├── indicators/indicators.go   # Technical indicators  
    ├── patterns/patterns.go       # Pattern detection  
    ├── dataloader/dataloader.go   # Data loading  
    ├── analyzer/analyzer.go       # Analysis engine  
    └── reporter/reporter.go       # **Report generation  

## ✨ Features Overview  
**Technical Indicators**  
### **RSI (Relative Strength Index)**  
Purpose: Momentum oscillator measuring price change velocity  
Range: 0-100 scale  
Signals:  
Above 70: Potentially overbought (sell signal)  
Below 30: Potentially oversold (buy signal)  
50: Neutral momentum  
Implementation: 14-period default with configurable timeframe  
Formula: RSI = 100 - (100 / (1 + RS)), where RS = Average Gain / Average Loss  
### **MACD (Moving Average Convergence Divergence)**  
Components:  
MACD Line: 12-period EMA - 26-period EMA  
Signal Line: 9-period EMA of MACD line  
Histogram: MACD line - Signal line  
Signals:  
MACD above signal: Bullish momentum  
MACD below signal: Bearish momentum  
Zero line crossover: Trend change confirmation  
Divergence: Potential reversal warning  
Advanced Features: Divergence detection, momentum strength analysis  
### **Bollinger Bands**  
Structure:  
Middle Band: 20-period Simple Moving Average  
Upper Band: Middle Band + (2 × Standard Deviation)  
Lower Band: Middle Band - (2 × Standard Deviation)  
Signals:  
Price touching upper band: Potential resistance/overbought  
Price touching lower band: Potential support/oversold  
Band squeeze: Low volatility, potential breakout  
Band expansion: High volatility period  
Analysis: Volatility measurement, mean reversion identification  
### **Moving Averages**  
**Simple Moving Average (SMA):**  
Arithmetic mean of closing prices  
Smooths price data to identify trends  
Less responsive to recent price changes  
**Exponential Moving Average (EMA):**  
Weighted average giving more importance to recent prices  
More responsive to price movements  
Better for short-term trend identification  
Applications: Trend direction, support/resistance levels, crossover signals  
### **Stochastic Oscillator**  
Components:  
%K Line: (Current Close - Lowest Low) / (Highest High - Lowest Low) × 100  
%D Line: 3-period SMA of %K  
Range: 0-100 scale  
Signals:  
Above 80: Overbought condition  
Below 20: Oversold condition  
%K crossing %D: Momentum change  
Analysis: Momentum measurement, reversal identification  
## Risk Assessment Metrics  
### Volatility Analysis  
**Anualized Volatility:**  
Standard deviation of daily returns × √252  
Measures price uncertainty over time  
Higher values indicate greater risk/opportunity  
**Historical Volatility:**  
Rolling volatility calculations  
Identifies volatility clusters  
Compares current vs historical volatility  
**Volatility Regime Detection:**  
Low volatility: Consolidation periods  
High volatility: Trending or news-driven periods  
## Risk-Adjusted Performance  
Sharpe Ratio:  
- Formula: (Portfolio Return - Risk-free Rate) / Portfolio Standard Deviation  
- Measures excess return per unit of risk  
- Values > 1.0 considered good, > 2.0 excellent  
Sortino Ratio:  
- Uses downside deviation instead of total volatility  
- Focus on harmful volatility only  
- Better measure for asymmetric return distributions  
Information Ratio:  
- Active return divided by tracking error  
- Measures risk-adjusted active return  
## Drawdown Analysis  
Maximum Drawdown:  
- Largest peak-to-trough decline  
- Worst-case scenario measurement  
- Critical for position sizing  
Average Drawdown:  
- Mean of all drawdown periods  
- Typical loss expectation  
Recovery Time Analysis:  
- Time to recover from drawdowns  
- Persistence of losses measurement  
Drawdown Duration:  
- Length of underwater periods  
- Risk tolerance assessment  
## Value at Risk (VaR) & Conditional VaR  
**95% VaR:**  
Maximum expected loss at 95% confidence  
Daily and monthly calculations  
Parametric and historical methods  
**99% VaR:**  
Extreme loss scenarios  
Stress testing measure  
**Conditional VaR (CVaR):**  
Expected loss beyond VaR threshold  
Tail risk measurement  
Expected shortfall calculation  
Implementation: Historical simulation, Monte Carlo methods  
# Pattern Recognition & Technical Analysis   
Dynamic Level Detection:  
Automatic identification of key price levels  
Historical price reaction analysis  
Strength scoring based on touches and volume  
Support Levels:  
Price floors where buying interest emerges  
Previous lows and consolidation areas  
Psychological round numbers  
Resistance Levels:  
Price ceilings where selling pressure appears  
Previous highs and supply zones  
Fibonacci levels and moving averages  
Level Validation:  
Multiple touches increase significance  
Volume confirmation at levels  
Time-based strength decay  
## Trend Analysis  
Trend Direction Detection:  
Algorithmic trend identification  
Multiple timeframe analysis  
Strength measurement (strong/weak/sideways)  
Uptrend Characteristics:  
Higher highs and higher lows  
Rising moving averages  
Positive momentum indicators  
Downtrend Characteristics:  
Lower highs and lower lows  
Declining moving averages  
Negative momentum indicators  
Sideways/Consolidation:  
Horizontal price movement  
Range-bound trading  
Low directional momentum  
## Candlestick Pattern Recognition  
**Single Candle Patterns:**    
Doji: Indecision, potential reversal  
Hammer: Bullish reversal at support  
Shooting Star: Bearish reversal at resistance  
Long White/Black: Strong directional movement  
**Two Candle Patterns:**    
Bullish Engulfing: Strong bullish reversal  
Bearish Engulfing: Strong bearish reversal  
Harami: Potential trend weakening  
**Three Candle Patterns:**    
Morning Star: Bullish reversal pattern  
Evening Star: Bearish reversal pattern  
Three White Soldiers: Strong bullish continuation  
## Volume Pattern Analysis  
Volume Breakout Detection:  
High volume with price movement  
Confirmation of trend changes  
Breakout validation  
Volume Divergence:  
Price vs volume relationship analysis  
Hidden strength/weakness detection  
Early warning signals  
Volume Patterns:  
Accumulation: Increasing volume with stable/rising prices  
Distribution: Increasing volume with stable/falling prices  
Low Volume: Lack of interest, potential reversal  
## Fibonacci Analysis  
Retracement Levels:  
23.6%, 38.2%, 50%, 61.8%, 76.4% levels  
Natural price correction points  
Support/resistance identification  
Extension Levels:  
127.2%, 161.8%, 261.8% projections  
Price target calculation  
Breakout level estimation  
Time-based Fibonacci:  
Fibonacci time zones  
Cycle analysis  
Timing reversal points  
## Data Source Integration  
### CoinGecko API Integration  
Real-time Data Access:  
Live Bitcoin market data  
Global price aggregation  
Multiple currency support (USD, EUR, BTC)  
Historical Data:  
Up to 365 days of history  
Hourly and daily granularity  
OHLCV data structure  
Rate Limiting:  
Respectful API usage  
Built-in retry mechanisms  
Error handling and fallbacks  
Data Quality:  
Professional-grade market data  
Multiple exchange aggregation  
Volume-weighted pricing  
### CSV/Excel Data Import  
Flexible Format Support:  
Auto-detection of column structure  
Multiple date/time formats  
Header row identification  
Supported Formats:  
OHLCV (Open, High, Low, Close, Volume)  
Timestamp-Price-Volume  
Extended formats with additional fields  
Data Validation:  
Missing data detection  
Outlier identification  
Chronological ordering verification  
Error Handling:  
Detailed error reporting  
Data quality warnings  
Suggestion for fixes  
### JSON Data Processing  
Structured Data Handling:  
Native JSON format support  
Nested object processing  
Array data extraction  
Schema Flexibility:  
Multiple JSON structures supported  
Custom field mapping  
Automatic type conversion  
Validation:  
Schema validation  
Required field checking  
Data type verification  
### Sample Data Generation  
Realistic Market Simulation:  
Mathematically generated price movements  
Geometric Brownian Motion model  
Volatility clustering simulation  
Configurable Parameters:  
Adjustable volatility levels  
Trend direction control  
Volume pattern simulation  
Testing Capabilities:  
Perfect for algorithm testing  
Consistent reproducible data  
Edge case scenario generation  
## Report Generation & Visualization  
### Interactive HTML Reports  
Modern Web Interface:  
Responsive design for all devices  
Professional styling and layout  
Print-friendly formatting  
Visual Elements:  
Color-coded trading signals  
Statistical summary tables  
Indicator visualization  
Interactive Features:  
Collapsible sections  
Hover tooltips  
Mobile-optimized navigation  
Technical Content:  
Complete analysis breakdown  
Signal explanations  
Risk metric interpretations  
### Machine-Readable JSON Output  
Structured Data:  
Complete analysis results  
Nested data organization  
API-ready format  
Integration Ready:  
Easy parsing for other systems  
Database storage compatible  
REST API integration  
Comprehensive Coverage:  
All calculated indicators  
Statistical measures  
Trading signals and reasoning  
### Console Output  
Quick Summary View:  
Key metrics at a glance  
Color-coded signals  
Progress indicators  
Verbose Mode:  
Detailed calculation steps  
Debug information  
Performance metrics  
Error Reporting:  
Clear error messages  
Troubleshooting guidance  
Data quality warnings  
## 🎛 Command Line Options  


`USAGE: btc-analyzer [OPTIONS]  

DATA SOURCE:  
  -source string    Data source: 'api', 'csv', 'json', 'sample' (default "api  
  -days int         Days for API data (default 30)  
  -csv string       CSV file path  
  -json string      JSON file path  

OUTPUT:  
  -output string    Output directory (default "output")  
  -html            Generate HTML report (default true)  
  -json-report     Generate JSON report (default true)  
  -verbose         Show detailed output  

EXAMPLES:  
  btc-analyzer -source=api -days=30  
  btc-analyzer -source=sample -days=60 -verbose  
  btc-analyzer -source=csv -csv=./data/prices.csv`  
## 📄 License  
MIT License - see LICENSE file for details.  

Disclaimer: This tool is for educational purposes only. Not financial advice. Always do your own research before making investment decisions.  


