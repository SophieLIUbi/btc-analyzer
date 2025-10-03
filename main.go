package main

import (
	"btc-analyzer/internal/analyzer"
	"btc-analyzer/internal/types"
	"btc-analyzer/internal/dataloader"
	"btc-analyzer/internal/reporter"
	"flag"
	"fmt"
	"log"
)

func main() {
	// Command line flags
	var (
		source     = flag.String("source", "api", "Data source: 'api', 'csv', 'json', or 'sample'")
		days       = flag.Int("days", 30, "Number of days for API data")
		csvFile    = flag.String("csv", "", "CSV file path")
		jsonFile   = flag.String("json", "", "JSON file path")
		outputDir  = flag.String("output", ".", "Output directory for reports")
		htmlReport = flag.Bool("html", true, "Generate HTML report")
		jsonReport = flag.Bool("json-report", true, "Generate JSON report")
		verbose    = flag.Bool("verbose", false, "Verbose output")
	)
	flag.Parse()

	fmt.Println("🚀 Bitcoin Market Analyzer Starting...")

	// Load data based on source
	var bts *types.BTCTimeSeries
	var err error

	switch *source {
	case "api":
		fmt.Printf("📡 Fetching %d days of data from CoinGecko API...\n", *days)
		bts, err = dataloader.LoadFromCoinGecko(*days)
		if err != nil {
			log.Fatalf("Failed to load data from API: %v", err)
		}

	case "csv":
		if *csvFile == "" {
			log.Fatal("CSV file path required when using -source=csv")
		}
		fmt.Printf("📄 Loading data from CSV file: %s\n", *csvFile)
		bts, err = dataloader.LoadFromCSV(*csvFile)
		if err != nil {
			log.Fatalf("Failed to load CSV data: %v", err)
		}

	case "json":
		if *jsonFile == "" {
			log.Fatal("JSON file path required when using -source=json")
		}
		fmt.Printf("📄 Loading data from JSON file: %s\n", *jsonFile)
		bts, err = dataloader.LoadFromJSON(*jsonFile)
		if err != nil {
			log.Fatalf("Failed to load JSON data: %v", err)
		}

	case "sample":
		fmt.Println("🎲 Generating sample data for demonstration...")
		bts = dataloader.GenerateSampleData(*days, 50000.0)

	default:
		log.Fatalf("Invalid source: %s. Use 'api', 'csv', 'json', or 'sample'", *source)
	}

	if bts == nil {
		log.Fatal("Failed to load data")
	}

	// Validate data
	fmt.Println("🔍 Validating data...")
	issues := dataloader.ValidateData(bts)
	if len(issues) > 0 {
		fmt.Printf("⚠️  Data validation warnings:\n")
		for _, issue := range issues {
			fmt.Printf("  - %s\n", issue)
		}
	} else {
		fmt.Println("✅ Data validation passed")
	}

	// Perform analysis
	fmt.Println("📊 Performing comprehensive analysis...")
	analytics := analyzer.PerformComprehensiveAnalysis(bts)

	// Print summary to console
	reporter.PrintSummary(bts, analytics)

	// Generate reports
	if *htmlReport {
		htmlPath := fmt.Sprintf("%s/btc_analysis_report.html", *outputDir)
		fmt.Printf("📝 Generating HTML report: %s\n", htmlPath)
		if err := reporter.GenerateHTMLReport(bts, analytics, htmlPath); err != nil {
			log.Printf("Failed to generate HTML report: %v", err)
		} else {
			fmt.Printf("✅ HTML report generated successfully\n")
		}
	}

	if *jsonReport {
		jsonPath := fmt.Sprintf("%s/btc_analysis_report.json", *outputDir)
		fmt.Printf("📝 Generating JSON report: %s\n", jsonPath)
		if err := reporter.GenerateJSONReport(bts, analytics, jsonPath); err != nil {
			log.Printf("Failed to generate JSON report: %v", err)
		} else {
			fmt.Printf("✅ JSON report generated successfully\n")
		}
	}

	// Save processed data
	csvPath := fmt.Sprintf("%s/btc_data.csv", *outputDir)
	fmt.Printf("💾 Saving data to CSV: %s\n", csvPath)
	if err := dataloader.SaveToCSV(bts, csvPath); err != nil {
		log.Printf("Failed to save CSV: %v", err)
	}

	if *verbose {
		fmt.Println("\n" + analyzer.GenerateReport(bts, analytics))
	}

	fmt.Println("🎉 Analysis complete! Check the output directory for reports.")
}