<template>
  <b-table
    v-if="logs.length > 0"
    :data="logs"
    :default-sort="['date', 'asc']"
    :per-page="settings.logsPerPage"
    detailed
    paginated
    pagination-simple
  >
    <template slot-scope="props">
      <b-table-column
        v-if="settings.cols.id"
        field="datetime"
        label="timestamp"
      >{{ props.row.date }}</b-table-column>
      <b-table-column
        v-if="settings.cols.hostname"
        field="source"
        label="source"
      >{{ props.row.source }}</b-table-column>
      <b-table-column
        v-if="settings.cols.pid"
        field="level"
        label="level"
      >{{ props.row.level }}</b-table-column>
      <b-table-column
        v-if="settings.cols.message"
        field="message"
        label="message"
      >{{ props.row.msg }}</b-table-column>
    </template>

    <template slot="detail" slot-scope="props">
      <span v-html="formatFullMessage(props.row)"/>
    </template>
  </b-table>
  <div v-else class="notification has-text-centered">
    <b-icon icon="package-variant" size="is-large"></b-icon>
    <p>
      No log messages received&hellip; yet!
    </p>
  </div>
</template>

<script>
export default {
  name: "logs-display",
  props: ["logs", "settings"],
  data() {
    return {
      todayDateString: new Date().toLocaleDateString()
    };
  },
  methods: {
    formatFullMessage(obj) {
      return "<pre>" + JSON.stringify(obj, null, "  ") + "</pre>";
    },
    severityClass(severity) {
      if (severity === 0) {
        // Emergicy
        return "is-dark";
      } else if (severity < 4) {
        // Alert, Critical, Error
        return "is-danger";
      } else if (severity === 4) {
        // Warning
        return "is-warning";
      } else if (severity === 7) {
        // Debug
        return "is-info";
      }

      return "is-light";
    }
  }
};
</script>
