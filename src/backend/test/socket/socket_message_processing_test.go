package socket_test

import (
	"encoding/binary"
	"net"
	"os"
	"testing"
	"time"

	"github.com/frr-mad/frr-mad/src/backend/internal/configs"
	"github.com/frr-mad/frr-mad/src/backend/internal/socket"
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

// TestMessageProcessing tests the socket's ability to process different message types
// and return appropriate responses
func TestMessageProcessing(t *testing.T) {
	config := configs.SocketConfig{
		UnixSocketLocation: "/tmp",
		UnixSocketName:     "test-message-socket",
		SocketType:         "unix",
	}

	os.Remove("/tmp/test-message-socket")

	mockLoggerInstance, mockAnalyzerInstance, mockMetrics, parsedAnalyzerdata := getMockData()

	socketInstance := socket.NewSocket(config, mockMetrics, mockAnalyzerInstance.AnalysisResult, mockLoggerInstance, parsedAnalyzerdata)

	go func() {
		socketInstance.Start()
	}()

	time.Sleep(100 * time.Millisecond)

	t.Run("TestOSPFDatabaseCommand", func(t *testing.T) {
		request := &frrProto.Message{
			Service: "ospf",
			Command: "database",
		}

		response := sendRequestAndGetResponse(t, request, "/tmp/test-message-socket")

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Returning OSPF database", response.Message)

		ospfDatabase := response.Data.GetOspfDatabase()
		assert.NotNil(t, ospfDatabase)
		assert.Equal(t, 1, len(ospfDatabase.Areas))
		assert.Equal(t, "192.168.1.1", ospfDatabase.Areas["0.0.0.0"].RouterLinkStates[0].Base.LsId)
	})

	t.Run("TestSystemResourcesCommand", func(t *testing.T) {
		request := &frrProto.Message{
			Service: "system",
			Command: "allResources",
		}

		response := sendRequestAndGetResponse(t, request, "/tmp/test-message-socket")

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Returning system metrics including CPU and memory", response.Message)

		systemMetrics := response.Data.GetSystemMetrics()
		assert.NotNil(t, systemMetrics)
		assert.Equal(t, 25.5, systemMetrics.CpuUsage)
		assert.Equal(t, 40.2, systemMetrics.MemoryUsage)
	})

	t.Run("TestInvalidCommand", func(t *testing.T) {
		request := &frrProto.Message{
			Service: "invalid",
			Command: "command",
		}

		response := sendRequestAndGetResponse(t, request, "/tmp/test-message-socket")

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "Unknown service: invalid")
	})

	t.Run("TestSystemDefault", func(t *testing.T) {
		request := &frrProto.Message{
			Service: "system",
			Command: "foo",
		}

		response := sendRequestAndGetResponse(t, request, "/tmp/test-message-socket")

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "There was an error getting system resources")

	})

	t.Run("TestSystemDefault", func(t *testing.T) {
		request := &frrProto.Message{
			Service: "system",
			Command: "exit",
		}

		response := sendRequestAndGetResponse(t, request, "/tmp/test-message-socket")

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Shutting system down", response.Message)

	})

	socketInstance.Close()
}

// Helper function to send a request to the socket and get the response
func sendRequestAndGetResponse(t *testing.T, request *frrProto.Message, socketPath string) *frrProto.Response {
	conn, err := net.Dial("unix", socketPath)
	assert.NoError(t, err)
	defer conn.Close()

	requestData, err := proto.Marshal(request)
	assert.NoError(t, err)

	sizeBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(sizeBuf, uint32(len(requestData)))

	_, err = conn.Write(sizeBuf)
	assert.NoError(t, err)
	_, err = conn.Write(requestData)
	assert.NoError(t, err)

	responseSizeBuf := make([]byte, 4)
	_, err = conn.Read(responseSizeBuf)
	assert.NoError(t, err)

	responseSize := binary.LittleEndian.Uint32(responseSizeBuf)
	responseData := make([]byte, responseSize)

	_, err = conn.Read(responseData)
	assert.NoError(t, err)

	response := &frrProto.Response{}
	err = proto.Unmarshal(responseData, response)
	assert.NoError(t, err)

	return response
}

func TestAnalysisHappyPath(t *testing.T) {
	s := getEmptyMockSocket()
	s.Anomalies.RouterAnomaly = CreateMockAnomalyDetectionRouter()
	s.Anomalies.ExternalAnomaly = CreateMockAnomalyDetectionExternal()
	s.Anomalies.NssaExternalAnomaly = CreateMockAnomalyDetectionNssaExternal()
	s.Anomalies.LsdbToRibAnomaly = CreateMockAnomalyDetectionLsdbToRib()
	m := &frrProto.Message{
		Service: "analysis",
		Command: "router",
	}

	t.Run("TestAnalysisRouter", func(t *testing.T) {
		m.Command = "router"
		response := s.ProcessCommand(m)
		assert.IsType(t, &frrProto.ResponseValue_Anomaly{}, response.Data.Kind)
		assert.Equal(t, "success", response.Status)
		assert.True(t, response.Data.GetAnomaly().HasOverAdvertisedPrefixes)
		assert.False(t, response.Data.GetAnomaly().HasUnAdvertisedPrefixes)
		assert.False(t, response.Data.GetAnomaly().HasDuplicatePrefixes)
		assert.False(t, response.Data.GetAnomaly().HasMisconfiguredPrefixes)
		assert.Equal(t, "89.207.132.170", response.Data.GetAnomaly().SuperfluousEntries[0].InterfaceAddress)
		assert.Empty(t, response.Data.GetAnomaly().MissingEntries)
		assert.Empty(t, response.Data.GetAnomaly().DuplicateEntries)
	})

	t.Run("TestAnalysisExternal", func(t *testing.T) {
		m.Command = "external"
		response := s.ProcessCommand(m)
		assert.IsType(t, &frrProto.ResponseValue_Anomaly{}, response.Data.Kind)
		assert.Equal(t, "success", response.Status)
		assert.False(t, response.Data.GetAnomaly().HasOverAdvertisedPrefixes)
		assert.True(t, response.Data.GetAnomaly().HasUnAdvertisedPrefixes)
		assert.False(t, response.Data.GetAnomaly().HasDuplicatePrefixes)
		assert.False(t, response.Data.GetAnomaly().HasMisconfiguredPrefixes)
		assert.Empty(t, response.Data.GetAnomaly().SuperfluousEntries)
		assert.Equal(t, "89.207.132.170", response.Data.GetAnomaly().MissingEntries[0].InterfaceAddress)
		assert.Empty(t, response.Data.GetAnomaly().DuplicateEntries)
	})

	t.Run("TestAnalysisNssaExternal", func(t *testing.T) {
		m.Command = "nssaExternal"
		response := s.ProcessCommand(m)
		assert.IsType(t, &frrProto.ResponseValue_Anomaly{}, response.Data.Kind)
		assert.Equal(t, "success", response.Status)
		assert.True(t, response.Data.GetAnomaly().HasOverAdvertisedPrefixes)
		assert.False(t, response.Data.GetAnomaly().HasUnAdvertisedPrefixes)
		assert.False(t, response.Data.GetAnomaly().HasDuplicatePrefixes)
		assert.False(t, response.Data.GetAnomaly().HasMisconfiguredPrefixes)
		assert.Equal(t, "89.207.132.170", response.Data.GetAnomaly().SuperfluousEntries[0].InterfaceAddress)
		assert.Empty(t, response.Data.GetAnomaly().MissingEntries)
		assert.Empty(t, response.Data.GetAnomaly().DuplicateEntries)
	})
	t.Run("TestAnalysisLsdbToRib", func(t *testing.T) {
		m.Command = "lsdbToRib"
		response := s.ProcessCommand(m)
		assert.IsType(t, &frrProto.ResponseValue_Anomaly{}, response.Data.Kind)
		assert.Equal(t, "success", response.Status)
		assert.False(t, response.Data.GetAnomaly().HasOverAdvertisedPrefixes)
		assert.False(t, response.Data.GetAnomaly().HasUnAdvertisedPrefixes)
		assert.False(t, response.Data.GetAnomaly().HasDuplicatePrefixes)
		assert.False(t, response.Data.GetAnomaly().HasMisconfiguredPrefixes)
		assert.Empty(t, response.Data.GetAnomaly().SuperfluousEntries)
		assert.Empty(t, response.Data.GetAnomaly().MissingEntries)
		assert.Empty(t, response.Data.GetAnomaly().DuplicateEntries)
	})

}

func TestAnalysisUnhappyPath(t *testing.T) {
	s := getEmptyMockSocket()
	m := &frrProto.Message{
		Service: "analysis",
		Command: "router",
	}

	t.Run("TestUnknownService", func(t *testing.T) {
		m.Command = "foobar"
		response := s.ProcessCommand(m)

		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Unknown command: foobar", response.Message)

	})
	t.Run("TestAnalysisSwitchCase", func(t *testing.T) {
		m.Command = "router"
		response := s.ProcessCommand(m)
		assert.IsType(t, &frrProto.ResponseValue_Anomaly{}, response.Data.Kind)
		assert.Equal(t, "success", response.Status)

		m.Command = "external"
		response = s.ProcessCommand(m)
		assert.IsType(t, &frrProto.ResponseValue_Anomaly{}, response.Data.Kind)
		assert.Equal(t, "success", response.Status)

		m.Command = "nssaExternal"
		response = s.ProcessCommand(m)
		assert.IsType(t, &frrProto.ResponseValue_Anomaly{}, response.Data.Kind)
		assert.Equal(t, "success", response.Status)

		m.Command = "lsdbToRib"
		response = s.ProcessCommand(m)
		assert.IsType(t, &frrProto.ResponseValue_Anomaly{}, response.Data.Kind)
		assert.Equal(t, "success", response.Status)

		m.Command = "ribToFib"
		response = s.ProcessCommand(m)
		assert.IsType(t, &frrProto.ResponseValue_Anomaly{}, response.Data.Kind)
		assert.Equal(t, "success", response.Status)

		m.Command = "shouldParsedLsdb"
		response = s.ProcessCommand(m)
		assert.IsType(t, &frrProto.ResponseValue_ParsedAnalyzerData{}, response.Data.Kind)
		assert.Equal(t, "success", response.Status)
	})

}

func TestFrrHappyPath(t *testing.T) {
	s := getEmptyMockSocket()
	s.Metrics = CreateMockFullFRRData()
	m := &frrProto.Message{
		Service: "frr",
	}

	t.Run("TestUnknownCommand", func(t *testing.T) {
		m.Command = "foobar"
		response := s.ProcessCommand(m)

		assert.Equal(t, "error", response.Status)
		assert.Equal(t, "Unknown command: foobar", response.Message)
	})

	t.Run("TestCommand_routerData", func(t *testing.T) {
		m.Command = "routerData"
		response := s.ProcessCommand(m)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Returning FRR meta data of router itself", response.Message)
		assert.Equal(t, "r101", response.Data.GetFrrRouterData().RouterName)
		assert.Equal(t, "192.168.1.1", response.Data.GetFrrRouterData().OspfRouterId)
	})

	t.Run("TestCommand_rib", func(t *testing.T) {
		m.Command = "rib"
		response := s.ProcessCommand(m)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Returning all routes (RIB)", response.Message)
		assert.IsType(t, &frrProto.ResponseValue_RoutingInformationBase{}, response.Data.Kind)

	})

	t.Run("TestCommand_ribfibSummary", func(t *testing.T) {
		m.Command = "ribfibSummary"
		response := s.ProcessCommand(m)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Returning route summaries of RIB and FIB", response.Message)
		assert.IsType(t, &frrProto.ResponseValue_RibFibSummaryRoutes{}, response.Data.Kind)

	})

}
func TestFrrUnhappyPath(t *testing.T) {
	s := getEmptyMockSocket()
	m := &frrProto.Message{
		Service: "frr",
	}
	t.Run("TestCommand_routerData", func(t *testing.T) {
		m.Command = "routerData"
		response := s.ProcessCommand(m)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Returning FRR meta data of router itself", response.Message)
		assert.Equal(t, "r101", response.Data.GetFrrRouterData().RouterName)
		assert.Equal(t, "192.168.1.1", response.Data.GetFrrRouterData().OspfRouterId)
		assert.Empty(t, response.Data.GetFrrRouterData())
	})

	t.Run("TestCommand_rib", func(t *testing.T) {
		m.Command = "rib"
		response := s.ProcessCommand(m)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Returning all routes (RIB)", response.Message)
		assert.IsType(t, &frrProto.ResponseValue_RoutingInformationBase{}, response.Data.Kind)
		assert.Empty(t, response.Data.GetRoutingInformationBase())

	})

	t.Run("TestCommand_ribfibSummary", func(t *testing.T) {
		m.Command = "ribfibSummary"
		response := s.ProcessCommand(m)

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Returning route summaries of RIB and FIB", response.Message)
		assert.IsType(t, &frrProto.ResponseValue_RibFibSummaryRoutes{}, response.Data.Kind)
		assert.Empty(t, response.Data.GetRibFibSummaryRoutes())
	})

}
