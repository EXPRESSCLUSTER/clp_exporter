package collector

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	mirrorSubsystem = "mirror"
)

type mirrorCollector struct{}

func init() {
	registCollector(mirrorSubsystem, NewMirrorCollector)
}

func NewMirrorCollector() (Collector, error) {
	return &mirrorCollector{}, nil
}

func getLastLine(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	var lastLine string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lastLine = scanner.Text()
	}

	if err := scanner.Err(); err != nil {
		return "", err
	}

	return lastLine, nil
}

func (c *mirrorCollector) Update(ch chan<- prometheus.Metric) error {
	metricType := prometheus.GaugeValue

	// Get all mirror disk resources from clp.conf file
	clpcfget_md := exec.Command("clpcfget", "-e", "/root/resource/md")
	output, err := clpcfget_md.Output()
	if err != nil {
		log.Fatalf("Failed to execute clpcfget -e /root/resource/md: %v", err)
	}
	log.Printf("Mirror disk resources:\n%s", strings.TrimSpace(string(output)))

	// FIXME
	mds := strings.Split(strings.TrimSpace(string(output)), "\n")
	for _, md := range mds {
		if strings.TrimSpace(md) != "" {
			log.Printf("Mirror Disk Resource: %s", md)
		}
		clpcfget_nmppath := exec.Command("clpcfget", "-g", "/root/resource/md@"+md+"/parameters/nmppath")
		output, err := clpcfget_nmppath.Output()
		if err != nil {
			log.Fatalf("Failed to create command for clpcfget -g /root/resource/md@%s/parameters/nmppath: %v", md, err)
			continue
		}

		// Change to lowercase (NMP -> nmp)
		deviceName := path.Base(string(output))
		deviceNameLower := strings.ToLower(deviceName)
		fmt.Printf("NMP :%s", deviceNameLower)
		filepath := fmt.Sprintf("/opt/nec/clusterpro/perf/disk/%s.cur", strings.TrimSpace(deviceNameLower))
		fmt.Printf("File Path: %s\n", filepath)
		lastLine, err := getLastLine(filepath)
		if err != nil {
			fmt.Printf("Failed to get last line: %v", err)
		}
		fmt.Printf("%s, Last Line: %s\n", strings.TrimSpace(deviceNameLower), lastLine)
		reader := csv.NewReader(strings.NewReader(lastLine))
		records, err := reader.Read()
		if err != nil {
			fmt.Printf("CSV parse error: %v", err)
			continue
		}

		//  1: Write, Total
		write_total := fmt.Sprintf("%s_write_total", md)
		writeTotal, err := strconv.ParseFloat(records[1], 64)
		if err != nil {
			fmt.Printf("Failed to parse write total value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, write_total),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, writeTotal,
		)

		//  2: Write, Avg
		write_avg := fmt.Sprintf("%s_write_avg", md)
		writeAvg, err := strconv.ParseFloat(records[2], 64)
		if err != nil {
			fmt.Printf("Failed to parse write avg value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, write_avg),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, writeAvg,
		)

		//  3: Read, Total
		read_total := fmt.Sprintf("%s_read_total", md)
		readTotal, err := strconv.ParseFloat(records[3], 64)
		if err != nil {
			fmt.Printf("Failed to parse read total value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, read_total),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, readTotal,
		)

		//  4: Read, Avg
		read_avg := fmt.Sprintf("%s_read_avg", md)
		readAvg, err := strconv.ParseFloat(records[4], 64)
		if err != nil {
			fmt.Printf("Failed to parse read avg value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, read_avg),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, readAvg,
		)

		//  5: Local Disk Write, Total
		local_disk_write_total := fmt.Sprintf("%s_local_disk_write_total", md)
		localDiskWriteTotal, err := strconv.ParseFloat(records[5], 64)
		if err != nil {
			fmt.Printf("Failed to parse local disk write total value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, local_disk_write_total),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, localDiskWriteTotal,
		)

		//  6: Local Disk Write, Avg
		local_disk_write_avg := fmt.Sprintf("%s_local_disk_write_avg", md)
		localDiskWriteAvg, err := strconv.ParseFloat(records[6], 64)
		if err != nil {
			fmt.Printf("Failed to parse local disk write avg value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, local_disk_write_avg),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, localDiskWriteAvg,
		)

		//  7: Local Disk Read, Total
		local_disk_read_total := fmt.Sprintf("%s_local_disk_read_total", md)
		localDiskReadTotal, err := strconv.ParseFloat(records[7], 64)
		if err != nil {
			fmt.Printf("Failed to parse local disk write avg value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, local_disk_read_total),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, localDiskReadTotal,
		)

		//  8: Local Disk Read, Avg
		local_disk_read_avg := fmt.Sprintf("%s_local_disk_read_avg", md)
		localDiskReadAvg, err := strconv.ParseFloat(records[8], 64)
		if err != nil {
			fmt.Printf("Failed to parse local disk write avg value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, local_disk_read_avg),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, localDiskReadAvg,
		)

		//  9: Send, Total
		send_total := fmt.Sprintf("%s_send_total", md)
		sendTotal, err := strconv.ParseFloat(records[9], 64)
		if err != nil {
			fmt.Printf("Failed to parse send total value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, send_total),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, sendTotal,
		)

		// 10: Send, Avg
		send_avg := fmt.Sprintf("%s_send_avg", md)
		sendAvg, err := strconv.ParseFloat(records[10], 64)
		if err != nil {
			fmt.Printf("Failed to parse send avg value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, send_avg),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, sendAvg,
		)

		// 11: Compress Ratio
		compress_ratio := fmt.Sprintf("%s_compress_ratio", md)
		compressRatio, err := strconv.ParseFloat(records[11], 64)
		if err != nil {
			fmt.Printf("Failed to parse compress ratio value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, compress_ratio),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, compressRatio,
		)

		// 12: Sync Time, Max
		sync_time_max := fmt.Sprintf("%s_sync_time_max", md)
		syncTimeMax, err := strconv.ParseFloat(records[12], 64)
		if err != nil {
			fmt.Printf("Failed to parse sync time max value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, sync_time_max),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, syncTimeMax,
		)

		// 13: Sync Time, Avg
		sync_time_avg := fmt.Sprintf("%s_sync_time_avg", md)
		syncTimeAvg, err := strconv.ParseFloat(records[13], 64)
		if err != nil {
			fmt.Printf("Failed to parse sync time avg value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, sync_time_avg),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, syncTimeAvg,
		)

		// 14: Sync Ack Time, Max
		sync_ack_time_max := fmt.Sprintf("%s_sync_ack_time_max", md)
		syncAckTimeMax, err := strconv.ParseFloat(records[14], 64)
		if err != nil {
			fmt.Printf("Failed to parse sync ack time max value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, sync_ack_time_max),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, syncAckTimeMax,
		)

		// 15: Sync Ack Time, Cur
		sync_ack_time_cur := fmt.Sprintf("%s_sync_ack_time_cur", md)
		syncAckTimeCur, err := strconv.ParseFloat(records[15], 64)
		if err != nil {
			fmt.Printf("Failed to parse sync ack time cur value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, sync_ack_time_cur),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, syncAckTimeCur,
		)

		// 16: Recovery Ack Time, Max
		recovery_ack_time_max := fmt.Sprintf("%s_recovery_ack_time_max", md)
		recoveryAckTimeMax, err := strconv.ParseFloat(records[16], 64)
		if err != nil {
			fmt.Printf("Failed to parse recovery ack time max value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, recovery_ack_time_max),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, recoveryAckTimeMax,
		)

		// 17: Recovery Ack Time, Max2
		recovery_ack_time_max2 := fmt.Sprintf("%s_recovery_ack_time_max2", md)
		recoveryAckTimeMax2, err := strconv.ParseFloat(records[17], 64)
		if err != nil {
			fmt.Printf("Failed to parse recovery ack time max2 value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, recovery_ack_time_max2),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, recoveryAckTimeMax2,
		)

		// 18: Recovery Ack Time Cur
		recovery_ack_time_cur := fmt.Sprintf("%s_recovery_ack_time_cur", md)
		recoveryAckTimeCur, err := strconv.ParseFloat(records[18], 64)
		if err != nil {
			fmt.Printf("Failed to parse recovery ack time cur value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, recovery_ack_time_cur),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, recoveryAckTimeCur,
		)

		// 19: SyncDiff, Max
		sync_diff_max := fmt.Sprintf("%s_sync_diff_max", md)
		syncDiffMax, err := strconv.ParseFloat(records[19], 64)
		if err != nil {
			fmt.Printf("Failed to parse sync diff max value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, sync_diff_max),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, syncDiffMax,
		)

		// 20: SyncDiff, Cur
		sync_diff_cur := fmt.Sprintf("%s_sync_diff_cur", md)
		syncDiffCur, err := strconv.ParseFloat(records[20], 64)
		if err != nil {
			fmt.Printf("Failed to parse sync diff cur value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, sync_diff_cur),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, syncDiffCur,
		)

		// 21: Send Queue, Max
		send_queue_max := fmt.Sprintf("%s_send_queue_max", md)
		sendQueueMax, err := strconv.ParseFloat(records[21], 64)
		if err != nil {
			fmt.Printf("Failed to parse send queue max value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, send_queue_max),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, sendQueueMax,
		)

		// 22: Send Queue, Max2
		send_queue_max2 := fmt.Sprintf("%s_send_queue_max2", md)
		sendQueueMax2, err := strconv.ParseFloat(records[22], 64)
		if err != nil {
			fmt.Printf("Failed to parse send queue max2 value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, send_queue_max2),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, sendQueueMax2,
		)

		// 23: Send Queue, Cur
		send_queue_cur := fmt.Sprintf("%s_send_queue_cur", md)
		sendQueueCur, err := strconv.ParseFloat(records[23], 64)
		if err != nil {
			fmt.Printf("Failed to parse send queue cur value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, send_queue_cur),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, sendQueueCur,
		)

		// 24: Request Queue, Max
		request_queue_max := fmt.Sprintf("%s_request_queue_max", md)
		requestQueueMax, err := strconv.ParseFloat(records[24], 64)
		if err != nil {
			fmt.Printf("Failed to parse request queue max value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, request_queue_max),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, requestQueueMax,
		)

		// 25: Request Queue, Max2
		request_queue_max2 := fmt.Sprintf("%s_request_queue_max2", md)
		requestQueueMax2, err := strconv.ParseFloat(records[25], 64)
		if err != nil {
			fmt.Printf("Failed to parse request queue max value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, request_queue_max2),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, requestQueueMax2,
		)

		// 26: Request Queue, Cur
		request_queue_cur := fmt.Sprintf("%s_request_queue_cur", md)
		requestQueueCur, err := strconv.ParseFloat(records[26], 64)
		if err != nil {
			fmt.Printf("Failed to parse request queue cur value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, request_queue_cur),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, requestQueueCur,
		)

		// 27: MDC HB Time, Max
		mdc_hb_time_max := fmt.Sprintf("%s_mdc_hb_time_max", md)
		mdcHbTimeMax, err := strconv.ParseFloat(records[27], 64)
		if err != nil {
			fmt.Printf("Failed to parse mdc hb time max value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, mdc_hb_time_max),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, mdcHbTimeMax,
		)

		// 28: MDC HB Time, Max2
		mdc_hb_time_max2 := fmt.Sprintf("%s_mdc_hb_time_max2", md)
		mdcHbTimeMax2, err := strconv.ParseFloat(records[28], 64)
		if err != nil {
			fmt.Printf("Failed to parse mdc hb time max2 value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, mdc_hb_time_max2),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, mdcHbTimeMax2,
		)

		// 29: MDC HB Time, Cur
		mdc_hb_time_cur := fmt.Sprintf("%s_mdc_hb_time_cur", md)
		mdcHbTimeCur, err := strconv.ParseFloat(records[29], 64)
		if err != nil {
			fmt.Printf("Failed to parse mdc hb time cur value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, mdc_hb_time_cur),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, mdcHbTimeCur,
		)

		// 30: Local-Write Waiting Recovery-Read Time, Total
		local_write_waiting_recovery_read_time_total := fmt.Sprintf("%s_local_write_waiting_recovery_read_time_total", md)
		localWriteWaitingRecoveryReadTimeTotal, err := strconv.ParseFloat(records[30], 64)
		if err != nil {
			fmt.Printf("Failed to parse local write waiting recovery read time total value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, local_write_waiting_recovery_read_time_total),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, localWriteWaitingRecoveryReadTimeTotal,
		)

		// 31: Local-Write Waiting Recovery-Read Time, Total2
		local_write_waiting_recovery_read_time_total2 := fmt.Sprintf("%s_local_write_waiting_recovery_read_time_total2", md)
		localWriteWaitingRecoveryReadTimeTotal2, err := strconv.ParseFloat(records[31], 64)
		if err != nil {
			fmt.Printf("Failed to parse local write waiting recovery read time total2 value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, local_write_waiting_recovery_read_time_total2),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, localWriteWaitingRecoveryReadTimeTotal2,
		)

		// 32: Recovery-Read Waiting Local-Write Time, Total
		recovery_read_waiting_local_write_time_total := fmt.Sprintf("%s_recovery_read_waiting_local_write_time_total", md)
		recoveryReadWaitingLocalWriteTimeTotal, err := strconv.ParseFloat(records[32], 64)
		if err != nil {
			fmt.Printf("Failed to parse recovery read waiting local write time total value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, recovery_read_waiting_local_write_time_total),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, recoveryReadWaitingLocalWriteTimeTotal,
		)

		// 33: Recovery-Read Waiting Local-Write Time, Total2
		recovery_read_waiting_local_write_time_total2 := fmt.Sprintf("%s_recovery_read_waiting_local_write_time_total2", md)
		recoveryReadWaitingLocalWriteTimeTotal2, err := strconv.ParseFloat(records[33], 64)
		if err != nil {
			fmt.Printf("Failed to parse recovery read waiting local write time total2 value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, recovery_read_waiting_local_write_time_total2),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, recoveryReadWaitingLocalWriteTimeTotal2,
		)

		// 34: Unmount Time, Max
		unmount_time_max := fmt.Sprintf("%s_unmount_time_max", md)
		unmountTimeMax, err := strconv.ParseFloat(records[34], 64)
		if err != nil {
			fmt.Printf("Failed to parse unmount time max value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, unmount_time_max),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, unmountTimeMax,
		)

		// 35: Unmount Time, Last
		unmount_time_last := fmt.Sprintf("%s_unmount_time_last", md)
		unmountTimeLast, err := strconv.ParseFloat(records[35], 64)
		if err != nil {
			fmt.Printf("Failed to parse unmount time last value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, unmount_time_last),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, unmountTimeLast,
		)

		// 36: Fsck Time, Max
		fsck_time_max := fmt.Sprintf("%s_fsck_time_max", md)
		fsckTimeMax, err := strconv.ParseFloat(records[36], 64)
		if err != nil {
			fmt.Printf("Failed to parse fsck time max value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, fsck_time_max),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, fsckTimeMax,
		)

		// 37: Fsck Time, Last
		fsck_time_last := fmt.Sprintf("%s_fsck_time_last", md)
		fsckTimeLast, err := strconv.ParseFloat(records[37], 64)
		if err != nil {
			fmt.Printf("Failed to parse fsck time last value: %v", err)
			continue
		}
		ch <- prometheus.MustNewConstMetric(
			prometheus.NewDesc(
				prometheus.BuildFQName(namespace, mirrorSubsystem, fsck_time_last),
				fmt.Sprintf("Mirror %s", md),
				nil, nil,
			),
			metricType, fsckTimeLast,
		)
	}

	return nil
}
