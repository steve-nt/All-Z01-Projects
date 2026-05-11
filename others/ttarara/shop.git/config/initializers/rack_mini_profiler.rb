# rack-mini-profiler shows the small "ms" timing panel on pages in development.
# It is optional. Disable by default; enable when you need performance debugging:
#
#   RACK_MINI_PROFILER=1 rails server
#
if defined?(Rack::MiniProfiler)
  Rack::MiniProfiler.config.enabled = ENV["RACK_MINI_PROFILER"].to_s == "1"
end
