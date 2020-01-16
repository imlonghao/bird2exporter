package bird

import (
	"regexp"
	"strings"
)

type parsed map[string]interface{}

var (
	regexNeighborAddress = regexp.MustCompile(`(?mi)^Neighbor address:\s+(.*)`)
	regexNeighborAs      = regexp.MustCompile(`(?mi)^Neighbor AS:\s+(.*)`)
	regexLocalAs         = regexp.MustCompile(`(?mi)^Local AS:\s+(.*)`)
	regexNeighborId      = regexp.MustCompile(`(?mi)^Neighbor ID:\s+(.*)`)
	regexSourceAddress   = regexp.MustCompile(`(?mi)^Source address:\s+(.*)`)
	regexChannel         = regexp.MustCompile(`(?mi)^Channel (.*)`)
	regexRoute           = regexp.MustCompile(`(?mi)^Routes:\s+(\d*) imported, (\d*) exported, (\d*) preferred`)
	regexImportUpdates   = regexp.MustCompile(`(?mi)^Import updates: +(\d+|---) +(\d+|---) +(\d+|---) +(\d+|---) +(\d+|---)`)
	regexImportWithdraws = regexp.MustCompile(`(?mi)^Import withdraws: +(\d+|---) +(\d+|---) +(\d+|---) +(\d+|---) +(\d+|---)`)
	regexExportUpdates   = regexp.MustCompile(`(?mi)^Export updates: +(\d+|---) +(\d+|---) +(\d+|---) +(\d+|---) +(\d+|---)`)
	regexExportWithdraws = regexp.MustCompile(`(?mi)^Export withdraws: +(\d+|---) +(\d+|---) +(\d+|---) +(\d+|---) +(\d+|---)`)
	regexNextHop         = regexp.MustCompile(`(?mi)^BGP Next hop:\s+(.*)`)
)

func (b *Bird) ShowProtocol() parsed {
	result := make(parsed, 0)
	b.write("show protocols")
	output := b.read()
	outputTable := strings.Split(output, "\n")[1:]
	for _, line := range outputTable {
		fields := strings.Fields(line)
		if len(fields) == 5 {
			fields = append(fields, "")
		}
		result[fields[0]] = parsed{
			"proto": fields[1],
			"table": fields[2],
			"state": fields[3],
			"since": fields[4],
			"info":  strings.Join(fields[5:], " "),
		}
	}
	return result
}

func (b *Bird) ShowProtocolAll() parsed {
	result := b.ShowProtocol()
	b.write("show protocols all")
	output := strings.Join(strings.Split(b.read(), "\n")[1:], "\n")
	outputTable := strings.Split(output, "\n\n")
	channels := make([]string, 0)
	for _, line := range outputTable {
		name := strings.Split(line, " ")[0]
		thisName := result[name].(parsed)
		if regexResult := regexNeighborAddress.FindStringSubmatch(line); len(regexResult) > 0 {
			thisName["neighbor_address"] = regexResult[1]
		}
		if regexResult := regexNeighborAs.FindStringSubmatch(line); len(regexResult) > 0 {
			thisName["neighbor_as"] = regexResult[1]
		}
		if regexResult := regexLocalAs.FindStringSubmatch(line); len(regexResult) > 0 {
			thisName["local_as"] = regexResult[1]
		}
		if regexResult := regexNeighborId.FindStringSubmatch(line); len(regexResult) > 0 {
			thisName["neighbor_id"] = regexResult[1]
		}
		if regexResult := regexSourceAddress.FindStringSubmatch(line); len(regexResult) > 0 {
			thisName["source_address"] = regexResult[1]
		}
		if regexResult := regexChannel.FindAllStringSubmatch(line, -1); len(regexResult) > 0 {
			for _, channel := range regexResult {
				channels = append(channels, channel[1])
			}
		}
		if regexResult := regexRoute.FindAllStringSubmatch(line, -1); len(regexResult) > 0 {
			for idx, result := range regexResult {
				thisChannel := make(parsed)
				if thisName[channels[idx]] != nil {
					thisChannel = thisName[channels[idx]].(parsed)
				}
				c := make(parsed)
				c["imported"] = result[1]
				c["exported"] = result[2]
				c["preferred"] = result[3]
				thisChannel["route"] = c
				thisName[channels[idx]] = thisChannel
			}
		}
		if regexResult := regexImportUpdates.FindAllStringSubmatch(line, -1); len(regexResult) > 0 {
			for idx, result := range regexResult {
				thisChannel := make(parsed)
				if thisName[channels[idx]] != nil {
					thisChannel = thisName[channels[idx]].(parsed)
				}
				c := make(parsed)
				c["received"] = result[1]
				c["rejected"] = result[2]
				c["filtered"] = result[3]
				c["ignored"] = result[4]
				c["accepted"] = result[5]
				thisChannel["import_updates"] = c
				thisName[channels[idx]] = thisChannel
			}
		}
		if regexResult := regexImportWithdraws.FindAllStringSubmatch(line, -1); len(regexResult) > 0 {
			for idx, result := range regexResult {
				thisChannel := make(parsed)
				if thisName[channels[idx]] != nil {
					thisChannel = thisName[channels[idx]].(parsed)
				}
				c := make(parsed)
				c["received"] = result[1]
				c["rejected"] = result[2]
				c["filtered"] = result[3]
				c["ignored"] = result[4]
				c["accepted"] = result[5]
				thisChannel["import_withdraws"] = c
				thisName[channels[idx]] = thisChannel
			}
		}
		if regexResult := regexExportUpdates.FindAllStringSubmatch(line, -1); len(regexResult) > 0 {
			for idx, result := range regexResult {
				thisChannel := make(parsed)
				if thisName[channels[idx]] != nil {
					thisChannel = thisName[channels[idx]].(parsed)
				}
				c := make(parsed)
				c["received"] = result[1]
				c["rejected"] = result[2]
				c["filtered"] = result[3]
				c["ignored"] = result[4]
				c["accepted"] = result[5]
				thisChannel["export_updates"] = c
				thisName[channels[idx]] = thisChannel
			}
		}
		if regexResult := regexExportWithdraws.FindAllStringSubmatch(line, -1); len(regexResult) > 0 {
			for idx, result := range regexResult {
				thisChannel := make(parsed)
				if thisName[channels[idx]] != nil {
					thisChannel = thisName[channels[idx]].(parsed)
				}
				c := make(parsed)
				c["received"] = result[1]
				c["rejected"] = result[2]
				c["filtered"] = result[3]
				c["ignored"] = result[4]
				c["accepted"] = result[5]
				thisChannel["export_withdraws"] = c
				thisName[channels[idx]] = thisChannel
			}
		}
		if regexResult := regexNextHop.FindStringSubmatch(line); len(regexResult) > 0 {
			thisName["next_hop"] = regexResult[1]
		}
		result[name] = thisName
	}
	return result
}
