export default function validate(value, regex, min, max, type) {
  const fail = () => ({ success: false });

  if (type) {
    const t = String(type).toLowerCase();
    if (t === 'int' || t === 'integer') {
      const n = Number(value);
      if (!Number.isInteger(n)) return fail();
      value = n;
    } else if (t === 'float' || t === 'number') {
      const n = Number(value);
      if (!isFinite(n)) return fail();
      value = n;
    } else if (t === 'string') {
      if (typeof value !== 'string') return fail();
    } else if (t === 'boolean' || t === 'bool') {
      if (typeof value !== 'boolean') {
        if (value === 'true') value = true;
        else if (value === 'false') value = false;
        else return fail();
      }
    } else if (t === 'array') {
      if (!Array.isArray(value)) return fail();
    } else if (t === 'object') {
      if (typeof value !== 'object' || value === null || Array.isArray(value)) return fail();
    }
  }

  if (regex) {
    let re = regex;
    if (typeof regex === 'string') {
      try { re = new RegExp(regex); } catch (e) { return fail(); }
    }
    if (!(re instanceof RegExp)) return fail();
    if (typeof value !== 'string' && typeof value !== 'number') return fail();
    if (!re.test(String(value))) return fail();
  }

  if (min !== undefined || max !== undefined) {
    if (typeof value === 'number') {
      if (min !== undefined && value < min) return fail();
      if (max !== undefined && value > max) return fail();
    } else if (typeof value === 'string' || Array.isArray(value)) {
      const len = value.length;
      if (min !== undefined && len < min) return fail();
      if (max !== undefined && len > max) return fail();
    } else if (typeof value === 'object' && value !== null) {
      const len = Object.keys(value).length;
      if (min !== undefined && len < min) return fail();
      if (max !== undefined && len > max) return fail();
    }
  }

  return { success: true };
}